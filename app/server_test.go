package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_application_routes(t *testing.T) {
	var registered = []struct {
		route  string
		method string
	}{
		{"/", "GET"},
		{"/static/*", "GET"},
		{"/login", "POST"},
	}

	var srv server
	mux := srv.routes()
	chiRoutes := mux.(chi.Routes)

	for _, route := range registered {
		if !routeExist(route.route, route.method, chiRoutes) {
			t.Errorf("route %s is not registered.", route.route)
		}
	}
}

func routeExist(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false
	_ = chi.Walk(chiRoutes, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		found = found || strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute)
		return nil
	})
	return found
}

func TestAppHome(t *testing.T) {
	tests := []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "", "From Session:"},
		{"second visit", "hello testing people", "From Session: hello testing people"},
	}
	for _, test := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		req = addContextAndSessionToRequest(req, &s)
		_ = s.Session.Destroy(req.Context())
		if test.putInSession != "" {
			s.Session.Put(req.Context(), "test", test.putInSession)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.Home)

		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("TestAppHome returned wrong status code. Expected %d, got %d", http.StatusOK, rr.Code)
		}
		body, _ := io.ReadAll(rr.Body)
		if !strings.Contains(string(body), test.expectedHTML) {
			t.Error("Did not find expected content in HTML")
		}
	}
}

func getContext(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, s *server) *http.Request {
	req = req.WithContext(getContext(req))
	c, _ := s.Session.Load(req.Context(), req.Header.Get("X-Session"))
	return req.WithContext(c)
}
