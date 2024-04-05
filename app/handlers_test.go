package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func Test_server_login(t *testing.T) {
	var tests = []struct {
		name               string
		postedData         url.Values
		expectedStatusCode int
		expectedLoc        string
	}{
		{
			name: "valid login",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/user/profile",
		},
		{
			name: "bad form submission",
			postedData: url.Values{
				"email":    {},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
		{
			name: "user not found",
			postedData: url.Values{
				"email":    {"admi@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
		{
			name: "bad password",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secreti"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(
			"POST",
			"/login",
			strings.NewReader(test.postedData.Encode()),
		)
		req = addContextAndSessionToRequest(req, &s)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedStatusCode {
			t.Errorf("%s: returned wrong status code; expected %d, but got %d",
				test.name, test.expectedStatusCode, rr.Code)
		}
		actualLoc, err := rr.Result().Location()
		if err != nil {
			t.Errorf("%s: no location header found: %v", test.name, err)
		} else {
			if actualLoc.Path != test.expectedLoc {
				t.Errorf("%s: returned wrong location; expected %s, but got %s",
					test.name, test.expectedLoc, actualLoc.String())
			}
		}
	}

}
