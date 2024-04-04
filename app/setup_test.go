package main

import (
	"os"
	"testing"
)

var s server

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
