package floop

import "testing"

func Test_Config(t *testing.T) {
	conf, err := LoadConfig("./config.yml")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", conf)
}
