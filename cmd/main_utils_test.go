package main

import (
	"bytes"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fithurriague/tax/internal/adapters/api/controller"
	"github.com/fithurriague/tax/internal/adapters/api/server"
)

func networkRequest(
	tcontroller controller.Controller,
	route controller.Route,
	payload []byte,
) *httptest.ResponseRecorder {
	srv := server.New(
		"",
		":8080",
		log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime),
		tcontroller,
	)

	endpoint := tcontroller.Endpoints()[route]

	r := httptest.NewRequest(string(endpoint.Method), string(route), bytes.NewReader(payload))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/json")

	srv.Handle(endpoint.Handler)(w, r)

	return w
}

func testNetworkRequest(
	tcontroller controller.Controller,
	route controller.Route,
	payload []byte,
	assert func(httptest.ResponseRecorder),
) func(t *testing.T) {
	return func(t *testing.T) {
		assert(*networkRequest(tcontroller, route, payload))
	}
}
