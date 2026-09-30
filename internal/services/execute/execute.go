package execute

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/merzzzl/cli-mcp-server/internal/models"
)

func (*Service) Execute(ctx context.Context, command models.Command) models.Result {
	if len(command.Arguments) == 0 {
		return models.Result{Text: "command arguments must not be empty", IsError: true}
	}

	if command.Timeout <= 0 {
		return models.Result{Text: "command timeout must be positive", IsError: true}
	}

	executionContext, cancel := context.WithTimeout(ctx, command.Timeout)
	defer cancel()

	process := exec.CommandContext(executionContext, command.Arguments[0], command.Arguments[1:]...)
	process.Dir = command.CWD

	output, executionErr := process.CombinedOutput()
	if executionErr == nil {
		return models.Result{Text: string(output)}
	}

	detail := fmt.Sprintf("%v", executionErr)
	if errors.Is(executionContext.Err(), context.DeadlineExceeded) {
		detail = fmt.Sprintf("command exceeded timeout of %s", command.Timeout)
	} else if errors.Is(executionContext.Err(), context.Canceled) {
		detail = "command was canceled"
	}

	text := string(output)
	if text != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	return models.Result{Text: text + detail, IsError: true}
}
