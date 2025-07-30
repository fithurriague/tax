package services

import (
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
	session := entities.MarketSession{}

	for _, op := range ops {
		// Buy
		if op.Type == entities.OperationTypeBuy {
			session.Buy(op.UnitCost, op.Quantity)
			continue
		}

		// Sell
		if op.Type == entities.OperationTypeSell {
			profit, err := session.Sell(op.UnitCost, op.Quantity)
			tax := profit * s.taxRate

			// Does NOT pay taxes
			if err != nil || profit <= 0 || op.Total() <= s.taxableAmount {
				taxes = append(taxes, entities.Tax{Tax: tax})
				return taxes, err
			}

			// Tax deduction
			if tax <= session.AccumulatedLosses {
				session.AccumulatedLosses -= tax
				tax = 0
			}

			taxes = append(taxes, entities.Tax{Tax: tax})
		}
	}

	return taxes, nil
}
