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

	"github.com/d3sw/floop/template"
	"github.com/d3sw/floop/types"
)

var (
	errInvalidURI    = "invalid uri: %s"
	errInvalidMethod = "invalid method: %d"
)

type endpointConfig struct {
	URI     string
	Method  string
	Body    string
	Headers map[string]string
}

// HTTPClientHandler implements a HTTP client handler for events
type HTTPClientHandler struct {
	conf   *endpointConfig
	client *http.Client
}

// NewHTTPClientHandler instantiates a new HTTPClientHandler
func NewHTTPClientHandler() *HTTPClientHandler {
	return &HTTPClientHandler{
		client: &http.Client{Timeout: 3 * time.Second},
	}
}

// Init initializes the http handler with the config
func (handler *HTTPClientHandler) Init(conf *types.HandlerConfig) error {
	config := conf.Config

	handler.conf = &endpointConfig{
		URI:     config["uri"].(string),
		Method:  config["method"].(string),
		Headers: make(map[string]string),
	}

	if _, ok := config["body"]; ok {
		handler.conf.Body = config["body"].(string)
	}

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

// Handle handles an event by making an http call per the config.
func (handler *HTTPClientHandler) Handle(event *types.Event) (map[string]interface{}, error) {
	resp, err := handler.httpDo(event, handler.conf)
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

func (handler *HTTPClientHandler) httpDo(event *types.Event, conf *endpointConfig) (*http.Response, error) {

	uri := template.Parse(event, conf.URI)
	body := template.Parse(event, conf.Body)
	buff := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(conf.Method, uri, buff)
	if err == nil {
		if conf.Headers != nil {
			for k, v := range conf.Headers {
				req.Header.Set(k, v)
			}
		}

		log.Printf("[DEBUG] handler=http uri='%s' body=%s", uri, body)
		return handler.client.Do(req)
	}

	return nil, err
}
