package handlers

import (
	"fmt"

	"github.com/d3sw/floop/types"
)

// EchoHandler implements a Handler that simply echoes back the input
type EchoHandler struct{}

// Handle echos back input data before process starts
func (lc *EchoHandler) Handle(event *types.Event) (map[string]interface{}, error) {
	fmt.Printf("[Echo] phase=%s %s\n", event.Type, event.Data)
	return nil, nil
}