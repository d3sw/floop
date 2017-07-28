package floop

import "testing"

func Test_Http_Config(t *testing.T) {
	conf, err := LoadConfig("./test-data/http.yml")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", conf)
}

func Test_Gnatsd_Config(t *testing.T) {
	conf, err := LoadConfig("./test-data/gnatsd.yml")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", conf)
}
