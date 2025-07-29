package controller

import (
	"encoding/json"
	"net/http"

	"github.com/fithurriague/tax/internal/domain/entities"
	"github.com/fithurriague/tax/internal/ports/iport"
)

type OperationController struct {
	key          string
	operationSvc iport.OperationService
}

func NewOperationController(
	operationSvc iport.OperationService,
) *OperationController {
	return &OperationController{
		key:          "operation_controller",
		operationSvc: operationSvc,
	}
}

func (c *OperationController) Key() string {
	return c.key
}

func (c *OperationController) Endpoints() map[Route]Endpoint {
	return map[Route]Endpoint{
		"/tax": {
			Method:  MethodPOST,
			Handler: c.Tax,
		},
	}
}

func (c *OperationController) Tax(
	w http.ResponseWriter,
	r *http.Request,
) Response {
	var ops []entities.Operation
	if err := json.NewDecoder(r.Body).Decode(&ops); err != nil {
		return Response{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	taxes, err := c.operationSvc.GetTaxes(ops)
	if err != nil {
		return Response{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	return Response{
		Status:  http.StatusOK,
		Content: taxes,
	}
}
