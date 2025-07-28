package controller

import (
	"net/http"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Endpoints() []Endpoint {
	return []Endpoint{
		{
			Route:   "/health",
			Method:  "GET",
			Handler: c.health,
		},
	}
}

func (c *HealthController) health(
	w http.ResponseWriter,
	r *http.Request,
) Response {
	return Response{
		Status: http.StatusOK,
		Content: "Up and running",
	}
}
