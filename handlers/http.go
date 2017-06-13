package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"github.com/d3sw/floop/lifecycle"
	"github.com/d3sw/floop/template"
)

var (
	errInvalidURI    = "invalid uri: %s"
	errInvalidMethod = "invalid method: %d"
)

// EndpointConfig is the config for a single endpoint
type EndpointConfig struct {
	URI    string
	Method string
	Body   string
}

// // HTTPConfig contains the http handler config
// type HTTPConfig struct {
// 	*EndpointConfig
// 	Begin     *EndpointConfig
// 	Progress  *EndpointConfig
// 	Failed    *EndpointConfig
// 	Completed *EndpointConfig
// }
//
// func (conf *HTTPConfig) setDefault(dst *EndpointConfig) {
// 	src := conf.EndpointConfig
//
// 	if src == nil || dst == nil {
// 		return
// 	}
//
// 	if dst.Method == "" {
// 		dst.Method = src.Method
// 	}
// 	if dst.URI == "" {
// 		dst.URI = src.URI
// 	}
// 	if dst.Body == "" {
// 		dst.Body = src.Body
// 	}
// }

// SetDefaults sets the defaults for values for config options not provided.
// func (conf *HTTPConfig) SetDefaults() {
// 	conf.setDefault(conf.Begin)
// 	conf.setDefault(conf.Progress)
// 	conf.setDefault(conf.Failed)
// 	conf.setDefault(conf.Completed)
// }

type HTTPClientHandler struct {
	conf   *EndpointConfig
	client *http.Client
}

// NewHTTPClientHandler instantiates a new HTTPClientHandler
func NewHTTPClientHandler(conf *EndpointConfig) *HTTPClientHandler {
	return &HTTPClientHandler{
		client: &http.Client{Timeout: 3 * time.Second},
		conf:   conf,
	}
}

func (handler *HTTPClientHandler) Handle(event *lifecycle.Event) error {
	resp, err := handler.httpDo(event, handler.conf)
	if err != nil {
		//log.Printf("[ERROR] %v", err)
		return err
	}

	if resp.StatusCode > 399 {
		//	log.Printf("[ERROR] %s", resp.Status)
		return errors.New(resp.Status)
	}

	return nil
}

func (handler *HTTPClientHandler) httpDo(event *lifecycle.Event, conf *EndpointConfig) (*http.Response, error) {
	//log.Println("Calling", conf)
	body := template.Parse(event, conf.Body)
	buff := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(conf.Method, conf.URI, buff)
	if err == nil {
		return handler.client.Do(req)
	}
	return nil, err
}

//
// // Progress makes a HTTP call with the event.  This is not intended to be used in the case of
// // realtime updates as it may be expensive.
// func (handler *HTTPClientHandler) Progress(event *lifecycle.Event) {
// 	conf := handler.conf.Progress
// 	if conf == nil {
// 		return
// 	}
//
// 	resp, err := handler.httpDo(event, conf)
// 	if err != nil {
// 		log.Printf("[ERROR] %v", err)
// 		return
// 	}
//
// 	if resp.StatusCode > 399 {
// 		log.Printf("[ERROR] %s", resp.Status)
// 		return
// 	}
// }
//
// func (handler *HTTPClientHandler) Completed(event *lifecycle.Event) {
// 	conf := handler.conf.Completed
// 	if conf == nil {
// 		return
// 	}
//
// 	resp, err := handler.httpDo(event, conf)
// 	if err != nil {
// 		log.Printf("[ERROR] %v", err)
// 		return
// 	}
//
// 	if resp.StatusCode > 399 {
// 		log.Printf("[ERROR] %s", resp.Status)
// 		return
// 	}
// }
//
// func (handler *HTTPClientHandler) Failed(event *lifecycle.Event) {
// 	conf := handler.conf.Failed
// 	if conf == nil {
// 		return
// 	}
//
// 	resp, err := handler.httpDo(event, conf)
// 	if err != nil {
// 		log.Printf("[ERROR] %v", err)
// 		return
// 	}
//
// 	if resp.StatusCode > 399 {
// 		log.Printf("[ERROR] %s", resp.Status)
// 		return
// 	}
// }
