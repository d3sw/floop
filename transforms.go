package floop

import "strings"

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
