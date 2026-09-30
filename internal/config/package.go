package config

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/merzzzl/cli-mcp-server/configs"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schema.json
var configSchema []byte

func Load(source string) (*Config, error) {
	var (
		data []byte
		err  error
	)

	if name, embedded := strings.CutPrefix(source, "@"); embedded {
		data, err = configs.Read(name)
	} else {
		data, err = os.ReadFile(source)
	}

	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", source, err)
	}

	schemaDocument, err := jsonschema.UnmarshalJSON(bytes.NewReader(configSchema))
	if err != nil {
		return nil, fmt.Errorf("decode config schema: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	if err = compiler.AddResource("config.schema.json", schemaDocument); err != nil {
		return nil, fmt.Errorf("add config schema: %w", err)
	}

	schema, err := compiler.Compile("config.schema.json")
	if err != nil {
		return nil, fmt.Errorf("compile config schema: %w", err)
	}

	document, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode config %q: %w", source, err)
	}

	if err = schema.Validate(document); err != nil {
		return nil, fmt.Errorf("validate config %q: %w", source, err)
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("decode validated config %q: %w", source, err)
	}

	for index := range config.Tools {
		tool := &config.Tools[index]
		cwd := map[string]any{
			"type":        "string",
			"minLength":   1,
			"description": "Working directory for the command",
		}

		switch properties := tool.Input["properties"].(type) {
		case nil:
			tool.Input["properties"] = map[string]any{"cwd": cwd}
		case map[string]any:
			properties["cwd"] = cwd
		default:
			continue
		}
	}

	return &config, nil
}
