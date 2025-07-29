package controller

import (
	"net/http"
)

type Method string

const (
	MethodGET     Method = "GET"
	MethodPOST    Method = "POST"
	MethodPUT     Method = "PUT"
	MethodPATCH   Method = "PATCH"
	MethodDELETE  Method = "DELETE"
	MethodOPTIONS Method = "OPTIONS"
)

type Response struct {
	Status  int   `json:"-"`
	Content any   `json:"content"`
	Err     error `json:"err,omitempty"`
}

type HandlerFunc func(http.ResponseWriter, *http.Request) Response

type Route string

type Endpoint struct {
	Method  Method
	Handler HandlerFunc
}

type Controller interface {
	Key() string
	Endpoints() map[Route]Endpoint
}
