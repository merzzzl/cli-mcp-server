package config

type Config struct {
	Name    string       `json:"name"`
	Version string       `json:"version"`
	Tools   []ConfigTool `json:"tools"`
}

type ConfigTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Input       map[string]any `json:"input"`
	Command     []string       `json:"command"`
	Timeout     string         `json:"timeout,omitempty"`
}
