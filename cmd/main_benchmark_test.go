package main

import (
	"encoding/json"
	"testing"

	"github.com/fithurriague/tax/internal/adapters/api/controller"
	"github.com/fithurriague/tax/internal/domain/entities"
	"github.com/fithurriague/tax/internal/domain/services"
)

func BenchmarkTax(b *testing.B) {
	buy1 := entities.Operation{
		Type:     entities.OperationTypeBuy,
		UnitCost: 10,
		Quantity: 1000,
	}

	buy2 := entities.Operation{
		Type:     entities.OperationTypeBuy,
		UnitCost: 20,
		Quantity: 2000,
	}

	sell1 := entities.Operation{
		Type:     entities.OperationTypeSell,
		UnitCost: 10,
		Quantity: 1000,
	}

	sell2 := entities.Operation{
		Type:     entities.OperationTypeSell,
		UnitCost: 30,
		Quantity: 2000,
	}

	noLossNoProfitOps := []entities.Operation{buy1, sell1}
	noLossNoProfitPayload, err := json.Marshal(noLossNoProfitOps)
	if err != nil {
		b.Fatal(err)
	}

	lossAndProfitOps := []entities.Operation{buy1, buy2, sell1, sell2}
	lossAndProfitPayload, err := json.Marshal(lossAndProfitOps)
	if err != nil {
		b.Fatal(err)
	}

	operationService := services.NewOperationService(entities.AllOperationTypes, 20000, 0.2)

	b.Run("No Loss No Profit", func(b *testing.B) {
		for b.Loop() {
			networkRequest(controller.NewOperationController(operationService), "/tax", noLossNoProfitPayload)
		}
	})

	b.Run("Loss And Profit", func(b *testing.B) {
		for b.Loop() {
			networkRequest(controller.NewOperationController(operationService), "/tax", lossAndProfitPayload)
		}
	})
}
