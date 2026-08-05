package router

import "github.com/shishir1290/goduck/runtime/handler"

type Route struct {
	Method  string
	Path    string
	Parts   []string
	Handler handler.HandlerFunc
}