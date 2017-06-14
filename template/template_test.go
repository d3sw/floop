package template

import (
	"testing"

	"github.com/d3sw/floop/types"
)

var testStr = `foo ${meta.key} ${type} something else ${data}`

func Test_Parse(t *testing.T) {

	event := &types.Event{
		Type: "begin",
		Meta: map[string]interface{}{"key": "value"},
		Data: []byte("foo"),
	}

	out := Parse(event, testStr)
	if out != "foo value begin something else foo" {
		t.Error("not parsed correctly")
	}
	event.Data = 5
	out = Parse(event, testStr)
	if out != "foo value begin something else 5" {
		t.Error("not parsed correctly")
	}
	t.Log(out)
}
