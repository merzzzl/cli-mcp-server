package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/merzzzl/cli-mcp-server/internal/config"
	"github.com/merzzzl/cli-mcp-server/internal/controller"
	"github.com/merzzzl/cli-mcp-server/internal/services/execute"
	"github.com/merzzzl/cli-mcp-server/internal/services/render"
)

func main() {
	arguments := os.Args[1:]
	if len(arguments) != 1 || arguments[0] == "" {
		log.Fatal("usage: cli-mcp-server <config.json|@name>")
	}

	configSource := arguments[0]
	if configSource[0] != '@' {
		absolutePath, err := filepath.Abs(configSource)
		if err != nil {
			log.Fatalf("resolve config path: %v", err)
		}

		configSource = absolutePath
	}

	configuration, err := config.Load(configSource)
	if err != nil {
		log.Fatal(err)
	}

	renderer := render.New(configuration)
	executor := execute.New()
	server := controller.New(configuration, renderer, executor)

	if err = server.Serve(context.Background()); err != nil {
		log.Fatal(err)
	}
}
