package shell

import (
	"fmt"
	"strings"

	"github.com/commander-cli/cmd"
	"github.com/fatih/color"

	"github.com/henomis/ai-shell-go/completion"
)

type Shell struct {
	completion *completion.Completion
}

type executeResponse string

const (
	executeResponseExecute executeResponse = "e"
	executeResponseRetry   executeResponse = "r"
	executeResponseExit    executeResponse = "q"
)

func New(completion *completion.Completion) *Shell {
	return &Shell{
		completion: completion,
	}
}

func (s *Shell) Run(input string) error {

	if input == "" {
		return fmt.Errorf("input is empty")
	}

	for {
		executeResponse, command := s.suggest(input)
		if executeResponse == executeResponseExit {
			c := cmd.NewCommand(command)
			err := c.Execute()
			if err != nil {
				return fmt.Errorf("command: %w", err)
			}
			break
		}

	}

	return nil
}

func (s *Shell) suggest(input string) (executeResponse, string) {
	response, err := s.completion.Suggest(input)
	if err != nil {
		return executeResponseExit, ""
	}

	color.New(color.FgCyan, color.Bold).Printf("$ %s\n", response.Command)
	color.New(color.FgWhite).Printf("--\n%s", response.Explain)

	var userInput string
	fmt.Println("[E]xecute, [R]etry, [Q]uit? > ")
	fmt.Scanf("%s", &userInput)

	switch strings.ToLower(userInput) {
	case "e":
		return executeResponseExecute, response.Command
	case "r":
		return executeResponseRetry, response.Command
	case "q":
		return executeResponseExit, response.Command
	default:
		return executeResponseExit, response.Command
	}

}
