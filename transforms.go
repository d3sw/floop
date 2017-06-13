package floop

import "strings"

func parsePerLineKeyValue(keyValuePairs, delim string) map[string]string {
	kvs := parseLines(keyValuePairs)
	kv := map[string]string{}
	// Parse key values delimited by '='
	for _, v := range kvs {
		arr := strings.Split(strings.TrimSpace(v), delim)
		if len(arr) > 1 {
			kv[strings.TrimSpace(arr[0])] = strings.TrimSpace(arr[1])
		}
	}

	return kv
}

func parseLines(line string) []string {
	lines := strings.Split(line, "\n")
	out := []string{}
	for _, line := range lines {
		if l := strings.TrimSpace(line); l != "" {
			out = append(out, line)
		}
	}
	return out
}

type transforms struct {
	KeyValue func(lines, delim string) map[string]string
	Line     func(lines string) []string
}

// Transforms holds an index of all available transforms
var Transforms = &transforms{
	KeyValue: parsePerLineKeyValue,
	Line:     parseLines,
}
