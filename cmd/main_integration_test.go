package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fithurriague/tax/internal/adapters/api/controller"
	"github.com/fithurriague/tax/internal/adapters/api/server"
	"github.com/fithurriague/tax/internal/domain/entities"
	"github.com/fithurriague/tax/internal/domain/services"
)

func TestTax(t *testing.T) {
	buy := entities.Operation{
		Type:     entities.OperationTypeBuy,
		UnitCost: 10,
		Quantity: 1000,
	}

	sell := entities.Operation{
		Type:     entities.OperationTypeSell,
		UnitCost: 10,
		Quantity: 1000,
	}

	ops := []entities.Operation{buy, sell}
	payload, err := json.Marshal(ops)
	if err != nil {
		t.Fatal(err)
	}

	operationService := services.NewOperationService(entities.AllOperationTypes, 20000, 0.2)

	srv := server.New(
		"/api",
		":8080",
		log.New(os.Stdout, "API: ", log.Ldate|log.Ltime),
		controller.NewHealthController(),
		controller.NewOperationController(operationService),
	)

	r := httptest.NewRequest(string(controller.MethodPOST), "/api/tax", bytes.NewReader(payload))
	w := httptest.NewRecorder()
	r.Header.Set("Content-Type", "application/json")

	endpoint := srv.Controllers["operation_controller"].Endpoints()["/tax"]
	srv.Handle(endpoint.Handler)(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
		return
	}

	var taxes []entities.Tax
	err = json.NewDecoder(w.Body).Decode(&taxes)
	if err != nil {
		t.Fatalf("Failed to decode response: %v, body: %s", err, w.Body.String())
	}

	t.Logf("RESPONSE: %v", taxes)
}
