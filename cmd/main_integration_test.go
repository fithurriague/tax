package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fithurriague/tax/internal/adapters/api/controller"
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

	t.Run(
		"TAX NETWORK REQUEST",
		testNetworkRequest(
			controller.NewOperationController(operationService),
			"/tax",
			payload,
			func(w httptest.ResponseRecorder) {
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
			},
		),
	)
}
