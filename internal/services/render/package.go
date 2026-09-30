package render

import (
	"errors"

	"github.com/merzzzl/cli-mcp-server/internal/config"
)

type Service struct {
	tools map[string]config.ConfigTool
}

var (
	ErrEmptyCommand   = errors.New("command is empty after rendering")
	ErrInvalidTimeout = errors.New("timeout must be a positive Go duration")
	ErrUnknownTool    = errors.New("unknown tool")
)

func New(configuration *config.Config) *Service {
	tools := make(map[string]config.ConfigTool, len(configuration.Tools))

	for toolIndex := range configuration.Tools {
		tool := configuration.Tools[toolIndex]
		tools[tool.Name] = tool
	}

	return &Service{tools: tools}
}
