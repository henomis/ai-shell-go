package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/commander-cli/cmd"
	"github.com/fatih/color"

	"github.com/henomis/ai-shell-go/completion"
)

type Shell struct {
	completion *completion.Completion
}

type ShellResponse struct {
	CommandAction CommandAction
	Command       string
}

type CommandAction string

const (
	CommandActionExecute CommandAction = "execute"
	CommandActionRevise  CommandAction = "retry"
	CommandActionExit    CommandAction = "exit"
)

func New(completion *completion.Completion) *Shell {
	return &Shell{
		completion: completion,
	}
}

func (s *Shell) Suggest(input string) (*ShellResponse, error) {

	if input == "" {
		return nil, fmt.Errorf("input is empty")
	}

	response, err := s.completion.Suggest(input, "")
	if err != nil {
		return nil, fmt.Errorf("completion: %w", err)
	}

	color.New(color.FgWhite, color.Bold).Printf("\n🤖 Here is your command line:\n\n")
	color.New(color.FgCyan, color.Bold).Printf("$ %s\n", response.Command)
	color.New(color.FgWhite, color.Italic).Printf("--\n%s\n\n", response.Explain)

	userAction := getUserActionFromStdin()

	return &ShellResponse{
		Command:       response.Command,
		CommandAction: getCommandActionFromUserAction(userAction),
	}, nil

}

func (s *Shell) Retry(previousCommand string) (*ShellResponse, error) {

	if previousCommand == "" {
		return nil, fmt.Errorf("command is empty")
	}
	var userAction string

	color.New(color.FgWhite, color.Bold).Printf("\n🤖 Enter your revision:\n\n")
	reader := bufio.NewReader(os.Stdin)
	userAction, _ = reader.ReadString('\n')

	response, err := s.completion.Suggest(userAction, previousCommand)
	if err != nil {
		return nil, fmt.Errorf("completion: %w", err)
	}

	color.New(color.FgWhite, color.Bold).Printf("\n🤖 Here is your command line:\n\n")
	color.New(color.FgCyan, color.Bold).Printf("$ %s\n", response.Command)
	color.New(color.FgWhite, color.Italic).Printf("--\n%s\n\n", response.Explain)

	userAction = getUserActionFromStdin()

	return &ShellResponse{
		Command:       response.Command,
		CommandAction: getCommandActionFromUserAction(userAction),
	}, nil

}

func (s *Shell) Execute(command string) error {

	color.New(color.FgCyan).DisableColor()

	shell := os.Getenv("SHELL")
	if shell != "" {
		command = fmt.Sprintf("%s -c '%s'", shell, command)
	}

	c := cmd.NewCommand(command, cmd.WithStandardStreams, cmd.WithInheritedEnvironment(cmd.EnvVars{}))

	err := c.Execute()
	if err != nil {
		return fmt.Errorf("command: %w", err)
	}

	return nil
}

// ---------------
// support methods
// ---------------

func getUserActionFromStdin() string {

	var userAction string
	color.New(color.FgWhite).Printf("[")
	color.New(color.FgGreen).Printf("E")
	color.New(color.FgWhite).Printf("]xecute, [")
	color.New(color.FgYellow).Printf("R")
	color.New(color.FgWhite).Printf("]evise, [")
	color.New(color.FgRed).Printf("Q")
	color.New(color.FgWhite).Printf("]uit? > ")
	reader := bufio.NewReader(os.Stdin)
	userAction, _ = reader.ReadString('\n')
	userAction = strings.TrimSpace(userAction)

	return userAction
}

func getCommandActionFromUserAction(userAction string) CommandAction {
	switch strings.ToLower(userAction) {
	case "e":
		return CommandActionExecute
	case "r":
		return CommandActionRevise
	case "q":
		return CommandActionExit
	default:
		return CommandActionExit
	}
}
