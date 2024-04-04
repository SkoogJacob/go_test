package main

import (
	"os"
	"testing"
)

var s server

func TestMain(m *testing.M) {
	s.Session = getSession()
	os.Exit(m.Run())
}
