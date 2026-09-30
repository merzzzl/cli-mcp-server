package render

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/merzzzl/cli-mcp-server/internal/models"
)

func (s *Service) Render(tool string, arguments map[string]any) (models.Command, error) {
	configuredTool, exists := s.tools[tool]
	if !exists {
		return models.Command{}, fmt.Errorf("%w %q", ErrUnknownTool, tool)
	}

	commandArguments := make([]string, 0, len(configuredTool.Command))
	for index, source := range configuredTool.Command {
		argumentTemplate, err := template.New(fmt.Sprintf("%s-%d", tool, index)).
			Option("missingkey=zero").
			Parse(source)
		if err != nil {
			return models.Command{}, fmt.Errorf("parse command argument %d: %w", index, err)
		}

		var rendered bytes.Buffer
		if err = argumentTemplate.Execute(&rendered, arguments); err != nil {
			return models.Command{}, fmt.Errorf("render command argument %d: %w", index, err)
		}

		if rendered.Len() != 0 {
			commandArguments = append(commandArguments, rendered.String())
		}
	}

	if len(commandArguments) == 0 {
		return models.Command{}, ErrEmptyCommand
	}

	cwd := ""
	if configuredCWD, configured := arguments["cwd"].(string); configured {
		cwd = configuredCWD
	}

	timeout := time.Minute

	if configuredTool.Timeout != "" {
		configuredTimeout, err := time.ParseDuration(configuredTool.Timeout)
		if err != nil {
			return models.Command{}, fmt.Errorf("parse command timeout: %w", err)
		}

		if configuredTimeout <= 0 {
			return models.Command{}, ErrInvalidTimeout
		}

		timeout = configuredTimeout
	}

	return models.Command{Arguments: commandArguments, CWD: cwd, Timeout: timeout}, nil
}
