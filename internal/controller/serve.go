package controller

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (c *Controller) Serve(ctx context.Context) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    c.configuration.Name,
		Version: c.configuration.Version,
	}, nil)

	for toolIndex := range c.configuration.Tools {
		tool := c.configuration.Tools[toolIndex]
		mcp.AddTool(server, &mcp.Tool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.Input,
		}, func(ctx context.Context, _ *mcp.CallToolRequest, arguments map[string]any) (*mcp.CallToolResult, any, error) {
			command, err := c.renderer.Render(tool.Name, arguments)
			if err != nil {
				return nil, nil, fmt.Errorf("render tool %s command: %w", tool.Name, err)
			}

			result := c.executor.Execute(ctx, command)

			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: result.Text}},
				IsError: result.IsError,
			}, nil, nil
		})
	}

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("serve MCP over stdio: %w", err)
	}

	return nil
}
