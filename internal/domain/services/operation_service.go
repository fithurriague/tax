package services

import (
	"errors"

	"github.com/fithurriague/tax/internal/domain/entities"
	"github.com/fithurriague/tax/internal/ports/iport"
)

type operationService struct {
	taxableOperations map[entities.OperationType]struct{}
	taxableAmount     float64
	taxRate           float64
}

func NewOperationService(
	taxableOperations []entities.OperationType,
	taxableAmount float64,
	taxRate float64,
) iport.OperationService {
	operationSvc := &operationService{
		taxableAmount:     taxableAmount,
		taxRate:           taxRate,
		taxableOperations: make(map[entities.OperationType]struct{}, len(taxableOperations)),
	}

	for _, op := range taxableOperations {
		operationSvc.taxableOperations[op] = struct{}{}
	}

	return operationSvc
}

func (s *operationService) GetTaxes(ops []entities.Operation) (taxes []entities.Tax, err error) {
	var accumulatedWeightedPrice float64 = 0
	buyedActions := 0
	var losses float64 = 0

	for _, op := range ops {
		// Buy
		if op.Type == entities.OperationTypeBuy {
			accumulatedWeightedPrice += op.UnitCost * float64(op.Quantity)
			buyedActions += op.Quantity
		}

		// Sell
		if op.Type == entities.OperationTypeSell {
			// Safe to assume it will never happen, but for correctness sake
			if op.Quantity > buyedActions {
				return nil, errors.New("not enough actions to sell")
			}

			averageUnitCost := accumulatedWeightedPrice / float64(buyedActions)
			if op.UnitCost <= averageUnitCost {
				// Accumulates losses
				diff := averageUnitCost - op.UnitCost
				losses += diff * float64(op.Quantity)
				// Add tax entry with 0 tax for loss or break-even operations
				taxes = append(taxes, entities.Tax{Tax: 0})
			} else {
				// Does not pay taxes
				if op.UnitCost*float64(op.Quantity) <= s.taxableAmount {
					taxes = append(taxes, entities.Tax{Tax: 0})
					continue
				}

				// Pays taxes
				diff := op.UnitCost - averageUnitCost
				profit := diff * float64(op.Quantity)
				tax := profit * s.taxRate

				// Tax deduction
				if tax <= losses {
					losses -= tax
					tax = 0
				}

				taxes = append(taxes, entities.Tax{Tax: tax})
			}
		}
	}

	return taxes, nil
}
