package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"web_test/pkg/data"
)

func Test_server_add_ip_to(t *testing.T) {
	tests := []struct {
		name        string
		headername  string
		headervalue string
		addr        string
		emptyAddr   bool
	}{
		{"default", "", "", "", false},
		{"default but with addr field", "", "", "", true},
		{"with header", "X-Forwarded-For", "192.0.1.42", "", false},
		{"spoof attempt", "", "", "spoofing:ongoing", false},
	}

	testHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(contextUserKey)
		if val == nil {
			t.Errorf("%s is not present", contextUserKey)
		}

		ip, ok := val.(string)
		if !ok {
			t.Error("value was not a string")
		}
		t.Log(ip)
	})

	for _, e := range tests {
		handlerToTest := s.addIpToContext(testHandlerFunc)

		req := httptest.NewRequest("GET", "http://testing.ex", nil)
		if e.emptyAddr {
			req.RemoteAddr = ""
		}

		if len(e.headername) > 0 {
			req.Header.Add(e.headername, e.headervalue)
		}
		if len(e.addr) > 0 {
			req.RemoteAddr = e.addr
		}

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func Test_server_ip_from_context(t *testing.T) {
	expected := "hello"
	c := context.WithValue(context.Background(), contextUserKey, expected)
	val := s.ipFromContext(c)
	if val != expected {
		t.Errorf("%s did not match the expected value %s", val, expected)
	}
}

func Test_server_auth(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	var tests = []struct {
		name   string
		isAuth bool
	}{
		{"logged in", true},
		{"not logged in", false},
	}

	for _, test := range tests {
		handlerToTest := s.auth(nextHandler)
		req := httptest.NewRequest("GET", "http://testing", nil)
		req = addContextAndSessionToRequest(req, &s)
		if test.isAuth {
			s.Session.Put(req.Context(), "user", data.User{ID: 1})
		}
		rr := httptest.NewRecorder()
		handlerToTest.ServeHTTP(rr, req)
		if test.isAuth && rr.Code != http.StatusOK {
			t.Errorf("%s: expected status code of 200 but got %d", test.name, rr.Code)
		} else if !test.isAuth && rr.Code != http.StatusSeeOther {
			t.Errorf("%s: expected status code is %d but got %d",
				test.name, http.StatusSeeOther, rr.Code)
		}
	}
}
