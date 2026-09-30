# CLI MCP Server

CLI MCP Server exposes command-line applications as MCP tools using declarative
JSON configurations. It communicates over stdio and starts processes directly,
without invoking a shell.

## Requirements

- Go 1.25 or newer;
- the CLI applications used by the selected configuration.

## Build

```sh
git clone https://github.com/merzzzl/cli-mcp-server.git
cd cli-mcp-server
go build -o cli-mcp-server ./cmd/server
```

## Run

The server requires exactly one configuration argument. It can be a filesystem
path:

```sh
./cli-mcp-server ./path/to/config.json
```

Or the name of an embedded configuration prefixed with `@`:

```sh
./cli-mcp-server @git
./cli-mcp-server @docker
./cli-mcp-server @kubectl
```

Embedded configurations are stored in `configs` and compiled into the binary.
For example, `@git` resolves to the embedded `configs/git.json` file.

## MCP client configuration

```json
{
  "mcpServers": {
    "git": {
      "command": "/absolute/path/to/cli-mcp-server",
      "args": ["@git"]
    }
  }
}
```

## Configuration

```json
{
  "name": "example-cli",
  "version": "1.0.0",
  "tools": [
    {
      "name": "example_status",
      "description": "Show application status",
      "input": {
        "type": "object",
        "properties": {
          "short": {
            "type": "boolean"
          }
        },
        "additionalProperties": false
      },
      "command": [
        "example",
        "status",
        "{{if .short}}--short{{end}}"
      ],
      "timeout": "10s"
    }
  ]
}
```

Top-level fields:

- `name` — MCP server name;
- `version` — MCP server version;
- `tools` — non-empty list of MCP tools.

Tool fields:

- `name` — unique MCP tool name;
- `description` — optional tool description;
- `input` — JSON Schema for the tool arguments;
- `command` — non-empty process argument list;
- `timeout` — optional positive Go duration, defaulting to one minute.

Every `command` item is an independent Go `text/template` and produces exactly
one process argument. Empty rendered items are omitted. Templates use
`missingkey=zero` and only the standard `text/template` functions are
available.

```text
{{.resource}}
{{with .name}}{{.}}{{end}}
{{if .short}}--short{{end}}
{{or .context "default-context"}}
```

The server automatically adds the optional string argument `cwd` to every tool
input. When provided, the child process runs in that directory. Configuration
files must not declare `cwd` themselves.

Standard output and standard error are returned together. A timeout or non-zero
exit code produces an MCP tool error.

## Security

Configuration files are trusted code because they determine which executables
MCP clients can start. Processes are launched without a shell, so shell control
operators in an argument are not interpreted by CLI MCP Server. The wrapped CLI
may still execute hooks, plugins, credential helpers, external tools, or other
programs according to its own configuration.
