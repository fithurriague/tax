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
	Status  int   `json:"status"`
	Content any   `json:"content"`
	Err     error `json:"err,omitempty"`
}

type HandlerFunc func(http.ResponseWriter, *http.Request) Response

type Endpoint struct {
	Route   string
	Method  Method
	Handler HandlerFunc
}

type Controller interface {
	Endpoints() []Endpoint
}
