package controller

import (
	"github.com/merzzzl/cli-mcp-server/internal/config"
	"github.com/merzzzl/cli-mcp-server/internal/services/execute"
	"github.com/merzzzl/cli-mcp-server/internal/services/render"
)

type Controller struct {
	configuration *config.Config
	renderer      *render.Service
	executor      *execute.Service
}

func New(
	configuration *config.Config,
	renderer *render.Service,
	executor *execute.Service,
) *Controller {
	return &Controller{
		configuration: configuration,
		renderer:      renderer,
		executor:      executor,
	}
}
