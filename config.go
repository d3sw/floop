package floop

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	"github.com/d3sw/floop/types"
)

// Config is the overall config.
type Config struct {
	Command       string
	Args          []string
	Meta          []string // required meta keys
	Quiet         bool
	ResolverHosts []string `yml:"resolverhost"`
	ResolverPort  int      `yml:"resolverport"`
	Handlers      map[types.EventType][]*types.HandlerConfig
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

// DefaultConfig returns a Config with defaults using the echo Lifecycle
func DefaultConfig() *Config {
	return &Config{
		Quiet:         false,
		ResolverHosts: []string{dResolverHost},
		ResolverPort:  dResolverPort,
		Handlers:      make(map[types.EventType][]*types.HandlerConfig),
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
