package floop

import (
	"fmt"
	"io/ioutil"

	"github.com/d3sw/floop/types"

	yaml "gopkg.in/yaml.v2"
)

// Config is the overall config.
type Config struct {
	Meta     []string // required meta keys
	Quiet    bool
	Handlers map[types.EventType][]*HandlerConfig
}

// HasMeta checks if the input meta has the required metadata keys
func (conf *Config) HasMeta(meta map[string]interface{}) bool {
	for _, m := range conf.Meta {
		if _, ok := meta[m]; !ok {
			return false
		}
	}
	return true
}

// TransformConfig holds configs for a transform
type TransformConfig []string

// HandlerConfig holds the config for a given handler
type HandlerConfig struct {
	Type      string                 // type of handler
	Transform TransformConfig        // transform to apply before making the callback
	Context   []string               // list of keys to set the context from callbacks
	Config    map[string]interface{} // Handler specific configs
	// Continue running child process even it handler returns error
	IgnoreErrors bool `yaml:"ignorerrors"`
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
	default:
		err = fmt.Errorf("unsupported transform: %s", conf.Transform[0])
	}

	return err
}

// DefaultConfig returns a Config with defaults using the echo Lifecycle
func DefaultConfig() *Config {
	return &Config{
		Quiet:    false,
		Handlers: make(map[types.EventType][]*HandlerConfig),
	}
}

// LoadConfig loads a config from the file given by filename
func LoadConfig(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var conf Config
	err = yaml.Unmarshal(b, &conf)
	return &conf, err
}
