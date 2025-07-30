package main

import (
	"encoding/json"
	"testing"

	"github.com/fithurriague/tax/internal/adapters/api/controller"
	"github.com/fithurriague/tax/internal/domain/entities"
	"github.com/fithurriague/tax/internal/domain/services"
)

func BenchmarkTax(b *testing.B) {
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
		b.Fatal(err)
	}

	operationService := services.NewOperationService(entities.AllOperationTypes, 20000, 0.2)

	for b.Loop() {
		networkRequest(controller.NewOperationController(operationService), "/tax", payload)
	}
}
