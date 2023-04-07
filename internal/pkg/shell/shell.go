package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/commander-cli/cmd"
	"github.com/fatih/color"

	"github.com/henomis/ai-shell-go/internal/pkg/completion"
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

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
)

func New(completion *completion.Completion) *Shell {
	return &Shell{
		completion: completion,
	}
}

func (s *Shell) Suggest(prompt string, previousCommand string) (*ShellResponse, error) {

	if prompt == "" {
		prompt = getUserPromptFromStdin()
	}

	return s.handleSuggestion(prompt, previousCommand)
}

func (s *Shell) Execute(command string) error {

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

func (s *Shell) handleSuggestion(prompt, previousCommand string) (*ShellResponse, error) {
	response, err := s.completion.Suggest(prompt, previousCommand)
	if err != nil {
		return nil, fmt.Errorf("completion: %w", err)
	}

	printCommandLineSuggestionFromResponse(response)

	return &ShellResponse{
		Command:       response.Command,
		CommandAction: newCommandActionFromUserAction(getUserActionFromStdin()),
	}, nil

}

// ---------------
// support methods
// ---------------

func printCommandLineSuggestionFromResponse(response *completion.CompletionResponse) {

	color.NoColor = false
	color.New(color.FgWhite, color.Bold).Printf("\n🤖 Here is your command line:\n\n")
	color.New(color.FgCyan, color.Bold).Printf("$ %s\n", response.Command)
	color.New(color.FgWhite, color.Italic).Printf("--\n%s\n\n", response.Explain)
	color.NoColor = true
}

func getUserPromptFromStdin() string {

	color.NoColor = false
	color.New(color.FgWhite, color.Bold).Printf("\n🤖 How may I help you? > ")
	color.NoColor = true

	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')

	return userInput
}

func getUserActionFromStdin() string {

	color.NoColor = false
	fmt.Printf("[%s]xecute, [%s]evise, [%s]uit? > ", green("E"), yellow("R"), red("Q"))
	color.NoColor = true

	reader := bufio.NewReader(os.Stdin)
	userAction, _ := reader.ReadString('\n')
	userAction = strings.TrimSpace(userAction)

	return userAction
}

func newCommandActionFromUserAction(userAction string) CommandAction {
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
