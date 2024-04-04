package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
)

type contextKey string

const contextUserKey contextKey = "user_ip"

func (s *server) ipFromContext(ctx context.Context) string {
	return ctx.Value(contextUserKey).(string)
}

func (s *server) addIpToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, err := getIp(r)
		if err != nil {
			log.Printf("Error encountered when parsing ip: %v", err)
		}
		ctx := context.WithValue(r.Context(), contextUserKey, ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getIp(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown", err
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "unknown", fmt.Errorf("unable to parse IP from %v", ip)
	}
	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}
	if len(ip) == 0 {
		ip = "forward"
	}
	return ip, nil
}
