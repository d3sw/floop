package types

import "fmt"

// TransformConfig holds configs for a transform
type TransformConfig []string

// HandlerConfig holds the config for a given handler
type HandlerConfig struct {
	// Type of handler
	Type string
	// Any uri
	URI string
	// Transform applied before calling the handler
	Transform TransformConfig
	// List of keys to set the context from handler responses
	Context []string
	// Body of the handler
	Body string
	// Handler specific configs
	Options map[string]interface{}
	// Continue running child process even it handler returns error
	IgnoreErrors bool `yaml:"ignorerrors"`
}

// Clone clones an existing config
func (conf *HandlerConfig) Clone() *HandlerConfig {
	return &HandlerConfig{
		Type:         conf.Type,
		URI:          conf.URI,
		Transform:    conf.Transform,
		Context:      conf.Context,
		Body:         conf.Body,
		Options:      conf.Options,
		IgnoreErrors: conf.IgnoreErrors,
	}
}

// ValidateTransform validates the transform
func (conf *HandlerConfig) ValidateTransform() error {
	if conf.Transform == nil || len(conf.Transform) == 0 {
		return nil
	}

	var err error
	switch conf.Transform[0] {
	case "kv":
		if len(conf.Transform) != 3 {
			err = fmt.Errorf("transform kv requires 2 arguments")
		}
	case "line":
		if len(conf.Transform) != 2 {
			err = fmt.Errorf("transform line requires 1 argument")
		}
	case "json":
		if len(conf.Transform) > 3 {
			err = fmt.Errorf("transform json invalid")
		}
	default:
		err = fmt.Errorf("transform unsupported: %s", conf.Transform[0])
	}

	return err
}
