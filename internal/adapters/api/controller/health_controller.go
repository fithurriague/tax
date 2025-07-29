package controller

import (
	"net/http"
)

type HealthController struct {
	key string
}

func NewHealthController() *HealthController {
	return &HealthController{
		key: "health_controller",
	}
}

func (c *HealthController) Key() string {
	return c.key
}

func (c *HealthController) Endpoints() map[Route]Endpoint {
	return map[Route]Endpoint{
		"/health": {
			Method:  MethodGET,
			Handler: c.health,
		},
	}
}

func (c *HealthController) health(
	w http.ResponseWriter,
	r *http.Request,
) Response {
	return Response{
		Status:  http.StatusOK,
		Content: "Up and running",
	}
}
