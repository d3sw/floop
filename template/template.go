package template

import (
	"fmt"
	"strings"

	"github.com/d3sw/floop/lifecycle"
)

func Parse(event *lifecycle.Event, str string) string {

	out := str
	for k, v := range event.Meta {
		old := "${meta." + k + "}"
		out = strings.Replace(out, old, v, -1)
	}

	out = strings.Replace(out, "${type}", string(event.Type), -1)

	var data string
	switch event.Data.(type) {
	case []byte:
		data = fmt.Sprintf("%s", event.Data)
	case int:
		data = fmt.Sprintf("%d", event.Data)
	}

	out = strings.Replace(out, "${data}", data, -1)

	return out
}
