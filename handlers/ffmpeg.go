package handlers

import (
	"log"
	"strings"

	"github.com/d3sw/floop/lifecycle"
)

// FFMPEGHandler implements a Handler for ffmpeg.  The begin and end events do nothing, only progress
// is used in this implementation
type FFMPEGHandler struct{}

// Handle handles progress callbacks for ffmpeg
func (handler *FFMPEGHandler) Handle(event *lifecycle.Event) error {
	b := event.Data.([]byte)
	kv := parsePerLineKeyValue(string(b), "=")

	if len(kv) > 0 {
		log.Printf("ffmpeg progress: %+v", kv)
	}

	return nil
}

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
