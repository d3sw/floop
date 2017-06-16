package floop

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/d3sw/floop/types"
)

var (
	errNoMatchingData       = errors.New("transform: no matching data")
	errUnsupportedTransform = errors.New("unsupported transform")
)

// Transform transforms the input given the transform and writes it to the event data
func Transform(transform []string, input []byte, out *types.Event) (transformed bool, err error) {
	switch transform[0] {
	case "kv":
		kvs := transformKeyValuePairs(string(input), transform[1], transform[2])
		if len(kvs) > 0 {
			out.Data = kvs
			transformed = true
		} else {
			err = errNoMatchingData
		}
	case "line":
		lines := transformLines(string(input), transform[1])
		if len(lines) > 0 {
			out.Data = lines
			transformed = true
		} else {
			err = errNoMatchingData
		}
	case "json":
		var v interface{}
		if err = json.Unmarshal(input, &v); err == nil {
			out.Data = v
			transformed = true
		}
	default:
		err = errUnsupportedTransform
	}

	return
}

// string to key-value map by pair and kv delimiter
func transformKeyValuePairs(keyValuePairs, kvpDelim, kvDelim string) map[string]string {
	kvs := transformLines(keyValuePairs, kvpDelim)
	kv := map[string]string{}
	// Parse key values delimited by delim
	for _, v := range kvs {
		arr := strings.Split(strings.TrimSpace(v), kvDelim)
		if len(arr) > 1 {
			kv[strings.TrimSpace(arr[0])] = strings.TrimSpace(arr[1])
		}
	}

	return kv
}

// string to string slice by delimiter
func transformLines(line string, delim string) []string {
	lines := strings.Split(line, delim)
	out := []string{}
	for _, line := range lines {
		if l := strings.TrimSpace(line); l != "" {
			out = append(out, line)
		}
	}
	return out
}
