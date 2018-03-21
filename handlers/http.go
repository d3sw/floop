package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/d3sw/floop/resolver"
	"github.com/d3sw/floop/types"
)

var (
	errInvalidURI    = "invalid uri: %s"
	errInvalidMethod = "invalid method: %d"
)

type endpointConfig struct {
	//URI     string
	Method string
	//Body    string
	Headers map[string]string
}

// Backoff interface defines contract for backoff strategies
type Backoff interface {
	Next(retry int) time.Duration
}

// LinearBackoff implements linear backoff
type LinearBackoff struct {
	Interval time.Duration
}

// Next returns next time for retrying operation with linear strategy
func (b LinearBackoff) Next(retry int) time.Duration {
	if retry <= 0 {
		return  time.Duration(0)
	}

	return  time.Duration(retry) * b.Interval
}

// HTTPClientHandler implements a HTTP client handler for events
type HTTPClientHandler struct {
	conf   *endpointConfig
	client *http.Client
	resolv *resolver.Resolver
	backoff Backoff
	retries int
}

// NewHTTPClientHandler instantiates a new HTTPClientHandler
func NewHTTPClientHandler(resolver *resolver.Resolver, backoff Backoff, maxRetries int) *HTTPClientHandler {
	return &HTTPClientHandler{
		client: &http.Client{Timeout: 3 * time.Second},
		resolv: resolver,
		backoff: backoff,
		retries: maxRetries,
	}
}

// Init initializes the http handler with the its specific config
func (handler *HTTPClientHandler) Init(conf *types.HandlerConfig) error {
	config := conf.Options

	handler.conf = &endpointConfig{
		//URI:     config["uri"].(string),
		//URI:     conf.URI,
		Method:  config["method"].(string),
		Headers: make(map[string]string),
	}

	//if _, ok := config["body"]; ok {
	//handler.conf.Body = config["body"].(string)
	//handler.conf.Body = string(conf.Body)
	//}

	if hdrs, ok := config["headers"]; ok {
		hm, ok := hdrs.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("invalid header data type %#v", config["headers"])
		}
		for k, v := range hm {
			key := k.(string)
			value := v.(string)
			handler.conf.Headers[key] = value
		}
	}

	return nil
}

// Handle handles an event by making an http call per the config.  Event is the raw event and
// HandlerConfig is the normalized config after interpolations have been applied.
func (handler *HTTPClientHandler) Handle(event *types.Event, conf *types.HandlerConfig) (map[string]interface{}, error) {
	resp, err := handler.httpDo(conf)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	defer resp.Body.Close()

	var r map[string]interface{}
	if err = json.Unmarshal(b, &r); err != nil {
		return nil, err
	}

	return r, nil
}

func (handler *HTTPClientHandler) httpDo(conf *types.HandlerConfig) (*http.Response, error) {
	buff := bytes.NewBuffer([]byte(conf.Body))

	var response *http.Response
	var err error

	attempt := 1
	for {
		discoveredURI, err := handler.resolv.Discover(conf.URI)
		if err != nil {
			log.Printf("[ERROR] Discovering URI [%s]: %s\n", conf.URI, err.Error())
			log.Println("[DEBUG] Will be used system DNS server")
		} else {
			conf.URI = discoveredURI
		}

		req, err := http.NewRequest(handler.conf.Method, conf.URI, buff)
		if err != nil {
			return nil, err
		}

		if handler.conf.Headers != nil {
			for k, v := range handler.conf.Headers {
				req.Header.Set(k, v)
			}
		}

		log.Printf("[DEBUG] handler=http uri='%s' body='%s'", conf.URI, conf.Body)
		response, err = handler.client.Do(req)
		if err == nil && response.StatusCode == 200{
			break
		}

		if attempt >= handler.retries {
			break
		}

		attempt ++
		time.Sleep(handler.backoff.Next(attempt))
	}

	return response, err
}

// CloseConnection - not implemented
func (handler *HTTPClientHandler) CloseConnection() error {
	//not implemented
	return nil
}
