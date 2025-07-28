package controller

import (
	"encoding/json"
	"net/http"

	"github.com/fithurriague/tax/internal/domain/entities"
)

type OperationController struct{}

func NewOperationController() *OperationController {
	return &OperationController{}
}

func (c *OperationController) Endpoints() []Endpoint {
	return []Endpoint{
		{
			Route:   "/tax",
			Method:  "POST",
			Handler: c.tax,
		},
	}
}

func (c *OperationController) tax(
	w http.ResponseWriter,
	r *http.Request,
) Response {
	var op entities.Operation
	if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
		return Response{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	// TODO: Perform operations

	return Response{
		Status: http.StatusOK,
		Content: "Hello from operation controller",
	}
}
