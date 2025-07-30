package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fithurriague/tax/internal/adapters/api/controller"
)

type Server struct {
	mux         *http.ServeMux
	logger      *log.Logger
	Address     string
	RoutePrefix string
	Controllers map[string]controller.Controller
}

func New(prefix string, addr string, logger *log.Logger, controllers ...controller.Controller) Server {
	srv := Server{
		mux:         http.NewServeMux(),
		logger:      logger,
		Address:     addr,
		RoutePrefix: prefix,
		Controllers: make(map[string]controller.Controller, len(controllers)),
	}

	for _, c := range controllers {
		srv.Controllers[c.Key()] = c
	}

	return srv
}

func (s *Server) Handle(hfn controller.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := hfn(w, r)

		// Set Content-Type header for JSON responses before WriteHeader
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.Status)

		var err error
		if contentType := w.Header().Get("Content-Type"); contentType == "text/plain" {
			if str, ok := res.Content.(string); ok {
				_, err = w.Write([]byte(str))
				if err == nil {
					return
				}
			}
		}

		// WARNING: BAD PRACTICE: Only done to comply with the challenge response format
		// I should ship the entire response as it is including content and error
		if res.Err != nil {
			err = json.NewEncoder(w).Encode(res.Err)
		} else {
			err = json.NewEncoder(w).Encode(res.Content)
		}

		if err != nil {
			s.logger.Printf("Failed to write response: %s", err.Error())
		}
	}
}

func (s *Server) ListenAndServe() {
	for _, cont := range s.Controllers {
		for route, endpoint := range cont.Endpoints() {
			fullRoute := fmt.Sprintf("%s %s%s", endpoint.Method, s.RoutePrefix, route)
			s.mux.Handle(fullRoute, s.Handle(endpoint.Handler))
		}
	}

	srv := &http.Server{
		Addr:    s.Address,
		Handler: s.mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		s.logger.Printf("Failed to start server: %v", err)
	}
}
