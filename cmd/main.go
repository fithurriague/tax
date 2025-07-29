package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fithurriague/tax/internal/adapters/api/controller"
	"github.com/fithurriague/tax/internal/adapters/api/server"
	"github.com/fithurriague/tax/internal/domain/entities"
	"github.com/fithurriague/tax/internal/domain/services"
)

func main() {
	operationService := services.NewOperationService(entities.AllOperationTypes, 20000, 0.2)

	srv := server.New(
		"/api",
		":8080",
		log.New(os.Stdout, "API: ", log.Ldate|log.Ltime),
		controller.NewHealthController(),
		controller.NewOperationController(operationService),
	)

	fmt.Printf("Starting server on: %s", srv.Address)
	srv.ListenAndServe()
}
