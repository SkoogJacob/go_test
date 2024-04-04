package main

import (
	"net/http"
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
