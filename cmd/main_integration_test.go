package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"

	"github.com/fithurriague/tax/internal/adapters/api/controller"
	"github.com/fithurriague/tax/internal/adapters/api/server"
	"github.com/fithurriague/tax/internal/domain/entities"
	"github.com/fithurriague/tax/internal/domain/services"
)

type TaxResponse struct {
	Status  int            `json:"-"`
	Content []entities.Tax `json:"content"`
	Err     error          `json:"err,omitempty"`
}

type TaxTestCase struct {
	Name  string
	Input []entities.Operation
	Want  TaxResponse
}

func TestTax(t *testing.T) {
	testcases := []TaxTestCase{
		{
			Name:  "EmptyOperations",
			Input: []entities.Operation{},
			Want: TaxResponse{
				Status:  http.StatusOK,
				Content: []entities.Tax{},
				Err:     nil,
			},
		},
		{
			Name: "ZeroUnitCost",
			Input: []entities.Operation{
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 0,
					Quantity: 100,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 10,
					Quantity: 100,
				},
			},
			Want: TaxResponse{
				Status:  http.StatusOK,
				Content: []entities.Tax{{Tax: 0}, {Tax: 0}},
				Err:     nil,
			},
		},
		{
			Name: "ZeroQuantity",
			Input: []entities.Operation{
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 10,
					Quantity: 0,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 15,
					Quantity: 0,
				},
			},
			Want: TaxResponse{
				Status:  http.StatusOK,
				Content: []entities.Tax{{Tax: 0}, {Tax: 0}},
				Err:     nil,
			},
		},
		{
			Name: "BuyOnly",
			Input: []entities.Operation{
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 10,
					Quantity: 100,
				},
			},
			Want: TaxResponse{
				Status:  http.StatusOK,
				Content: []entities.Tax{{Tax: 0}},
				Err:     nil,
			},
		},
		{
			Name: "BreakEven",
			Input: []entities.Operation{
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 10,
					Quantity: 100,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 10,
					Quantity: 100,
				},
			},
			Want: TaxResponse{
				Status:  http.StatusOK,
				Content: []entities.Tax{{Tax: 0}, {Tax: 0}},
				Err:     nil,
			},
		},
		{
			Name: "ProfitBelowTaxableAmount",
			Input: []entities.Operation{
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 10,
					Quantity: 100,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 15,
					Quantity: 10,
				},
			},
			Want: TaxResponse{
				Status:  http.StatusOK,
				Content: []entities.Tax{{Tax: 0}, {Tax: 0}},
				Err:     nil,
			},
		},
		{
			Name: "ProfitAboveTaxableAmount",
			Input: []entities.Operation{
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 10,
					Quantity: 2000,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 15,
					Quantity: 2000,
				},
			},
			Want: TaxResponse{
				Status:  http.StatusOK,
				Content: []entities.Tax{{Tax: 0}, {Tax: 2000}},
				Err:     nil,
			},
		},
		{
			Name: "Loss",
			Input: []entities.Operation{
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 15,
					Quantity: 100,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 10,
					Quantity: 100,
				},
			},
			Want: TaxResponse{
				Status:  http.StatusOK,
				Content: []entities.Tax{{Tax: 0}, {Tax: 0}},
				Err:     nil,
			},
		},
		{
			Name: "TaxDeductionBreakEven",
			Input: []entities.Operation{
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 15,
					Quantity: 100,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 10,
					Quantity: 100,
				},
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 10,
					Quantity: 100,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 15,
					Quantity: 100,
				},
			},
			Want: TaxResponse{
				Status:  http.StatusOK,
				Content: []entities.Tax{{Tax: 0}, {Tax: 0}, {Tax: 0}, {Tax: 0}},
				Err:     nil,
			},
		},
		{
			Name: "Complex",
			Input: []entities.Operation{
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 10,
					Quantity: 10000,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 2,
					Quantity: 5000,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 20,
					Quantity: 2000,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 20,
					Quantity: 2000,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 25,
					Quantity: 1000,
				},
				{
					Type:     entities.OperationTypeBuy,
					UnitCost: 20,
					Quantity: 10000,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 5,
					Quantity: 5000,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 30,
					Quantity: 4350,
				},
				{
					Type:     entities.OperationTypeSell,
					UnitCost: 30,
					Quantity: 650,
				},
			},
			Want: TaxResponse{
				Status:  http.StatusOK,
				Content: []entities.Tax{
					{Tax: 0},
					{Tax: 0},
					{Tax: 0},
					{Tax: 0},
					{Tax: 3000.0},
					{Tax: 0},
					{Tax: 0},
					{Tax: 3050.0},
					{Tax: 0},
				},
				Err:     nil,
			},
		},
	}

	const route = "/tax"
	operationService := services.NewOperationService(entities.AllOperationTypes, 20000, 0.2)
	operationController := controller.NewOperationController(operationService)
	endpoint := operationController.Endpoints()[route]

	srv := server.New(
		"",
		":8080",
		log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime),
		operationController,
	)

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			payload, err := json.Marshal(tc.Input)
			if err != nil {
				t.Fatal(err)
			}

			r := httptest.NewRequest(string(endpoint.Method), string(route), bytes.NewReader(payload))
			w := httptest.NewRecorder()
			r.Header.Set("Content-Type", "application/json")

			srv.Handle(endpoint.Handler)(w, r)

			if w.Code != tc.Want.Status {
				t.Errorf("Status: got = %v want %v", w.Code, tc.Want.Status)
			}

			resBody := w.Body.String()
			t.Logf("raw json response: %s", resBody)

			// Check for content
			var resContent []entities.Tax
			err = json.Unmarshal([]byte(resBody), &resContent)
			if err == nil {
				if !slices.Equal(resContent, tc.Want.Content) {
					t.Fatalf("Content: got = %v want = %v", resContent, tc.Want.Content)
				}

				return
			}

			// Check for error
			var resErr error
			err = json.Unmarshal([]byte(resBody), &resErr)
			if err == nil {
				if resErr != tc.Want.Err {
					t.Fatalf("Error: got = %v want = %v", resErr, tc.Want.Err)
				}
			}
		})
	}
}
