package template

import (
	"fmt"
	"strings"

	"github.com/d3sw/floop"
)

// Parse applies event context to the input string returning the normalized string.
func Parse(event *floop.Event, str string) string {

	out := str
	for k, v := range event.Meta {
		old := "${meta." + k + "}"
		var val string
		switch v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			val = fmt.Sprintf("%d", v)
		case float32, float64:
			val = fmt.Sprintf("%f", v)
		case []byte:
			bv := v.([]byte)
			val = string(bv)
		case string:
			val = v.(string)
		default:
			continue
		}
		out = strings.Replace(out, old, val, -1)
	}

	out = strings.Replace(out, "${type}", string(event.Type), -1)

	var data string
	switch event.Data.(type) {
	case []byte:
		bstr := event.Data.([]byte)
		data = string(bstr)
	case int:
		data = fmt.Sprintf("%d", event.Data)
	case string:
		data = event.Data.(string)
	}

	out = strings.Replace(out, "${data}", data, -1)

	return out
}
