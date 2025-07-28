package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fithurriague/tax/internal/adapters/api/controller"
	"github.com/fithurriague/tax/internal/adapters/api/server"
)

func main() {
	srv := server.New(
		"/api",
		":8080",
		log.New(os.Stdout, "API: ", log.Ldate|log.Ltime),
		controller.NewHealthController(),
		controller.NewOperationController(),
	)

	fmt.Printf("Starting server on: %s", srv.Address)
	srv.ListenAndServe()
}
