package floop

import (
	"io/ioutil"

	"github.com/d3sw/floop/lifecycle"

	yaml "gopkg.in/yaml.v2"
)

// Config is the overall config.
type Config struct {
	Quiet    bool
	Handlers map[lifecycle.EventType][]HandlerConfig
}

type HandlerConfig struct {
	Type      string
	Transform string
	Config    map[string]interface{}
}

// HandlersConfig holds the config for all handlers
// type HandlersConfig struct {
// 	HTTP   *handlers.HTTPConfig
// 	FFMPEG interface{}
// }

// DefaultConfig returns a Config with defaults using the echo Lifecycle
func DefaultConfig() *Config {
	return &Config{
		Quiet:    false,
		Handlers: make(map[lifecycle.EventType][]HandlerConfig),
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
