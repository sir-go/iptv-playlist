package main

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestE2e(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		dataPath string
	}{
		{"m3u-uc", "/pl.m3u", "../../testdata/pl.m3u"},
		{"m3u-mc", "/px.pl.m3u", "../../testdata/px.pl.m3u"},
		{"xspf-uc", "/pl.xspf", "../../testdata/pl.xspf"},
		{"xspf-mc", "/px.pl.xspf", "../../testdata/px.pl.xspf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.Get("http://localhost" + tt.url)
			if err != nil {
				t.Error(err)
			}

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Error(err)
			}

			dBytes, err := ioutil.ReadFile(tt.dataPath)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(dBytes, body) {
				t.Errorf("Got playlist (%v) not the same that in %v (%v)", len(body), tt.dataPath, len(dBytes))
			}

		})
	}
}
