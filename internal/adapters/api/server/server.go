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
	Controllers []controller.Controller
}

func New(prefix string, addr string, logger *log.Logger, controllers ...controller.Controller) Server {
	return Server{
		mux:         http.NewServeMux(),
		logger:      logger,
		Address:     addr,
		RoutePrefix: prefix,
		Controllers: controllers,
	}
}

func (s *Server) handle(hfn controller.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := hfn(w, r)

		var err error
		if contentType := w.Header().Get("Content-Type"); contentType == "text/plain" {
			if str, ok := res.Content.(string); ok {
				_, err = w.Write([]byte(str))
				if err == nil {
					return
				}
			}
		}

		err = json.NewEncoder(w).Encode(res.Content)
		if err != nil {
			s.logger.Printf("Failed to write response: %s", err.Error())
		}
	}
}

func (s *Server) ListenAndServe() {
	for _, cont := range s.Controllers {
		for _, endpoint := range cont.Endpoints() {
			route := fmt.Sprintf("%s %s%s", endpoint.Method, s.RoutePrefix, endpoint.Route)
			s.mux.Handle(route, s.handle(endpoint.Handler))
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
