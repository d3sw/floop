package floop

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config is the overall config.
type Config struct {
	Quiet    bool
	Handlers map[EventType][]*HandlerConfig
}

// HandlerConfig holds the config for a given handler
type HandlerConfig struct {
	Type    string                 // type of handler
	Context []string               // list of keys to set the context from callbacks
	Config  map[string]interface{} // Handler specific configs
}

// DefaultConfig returns a Config with defaults using the echo Lifecycle
func DefaultConfig() *Config {
	return &Config{
		Quiet:    false,
		Handlers: make(map[EventType][]*HandlerConfig),
	}
}

func LoadConfig(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var conf Config
	err = yaml.Unmarshal(b, &conf)
	return &conf, err
}
