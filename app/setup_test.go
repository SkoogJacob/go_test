package main

import (
	"encoding/gob"
	"os"
	"testing"
	"web_test/pkg/data"
	"web_test/pkg/repository/dbrepo"
)

var s server

func TestMain(m *testing.M) {
	pathToTemplates = "../templates/"
	// Set up DB connection

	s.DB = &dbrepo.TestRepo{
		Users: []*data.User{
			{ID: 1, Email: "admin@example.com"},
		},
	}

	s.Session = getSession()
	gob.Register(data.User{})

	os.Exit(m.Run())
}
