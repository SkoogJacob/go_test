package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var tests = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/fish", http.StatusNotFound},
	}
	routes := s.routes()
	pathToTemplates = "../templates/"

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, test := range tests {
		resp, err := ts.Client().Get(ts.URL + test.url)
		if err != nil {
			t.Log(err)
			t.Fatalf("Serve failed with error %v", err)
		}
		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("for %s expected status code %d but got %d", test.url, test.expectedStatusCode, resp.StatusCode)
		}
	}
}
