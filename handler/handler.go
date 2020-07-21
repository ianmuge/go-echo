package handler

import (
	"gopkg.in/mgo.v2"
	"os"
)

type (
	Handler struct {
		DB *mgo.Session
	}
)

var (
// Key (Should come from somewhere else).
	Key = os.Getenv("SECRET_KEY")
)