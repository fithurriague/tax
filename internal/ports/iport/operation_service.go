package iport

import "github.com/fithurriague/tax/internal/domain/entities"

type OperationService interface {
	GetTaxes(ops []entities.Operation) (taxes []entities.Tax, err error)
}
