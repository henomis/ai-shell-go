package main

import (
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"github.com/henomis/ai-shell-go/internal/pkg/completion"
	"github.com/henomis/ai-shell-go/internal/pkg/shell"
)

var (
	ErrorShellAI = fmt.Errorf("🤖 OOPS! ")
)

func main() {

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		fmt.Printf("%s: OPEN_AI_KEY is not set. Please set the OPENAI_API_KEY environment variable to your OpenAI API key\n", ErrorShellAI)
		return
	}

	userInput := strings.Join(os.Args[1:], " ")

	openAIClient := openai.NewClient(openAIKey)
	completionInstance := completion.New(openAIClient)
	shellInstance := shell.New(completionInstance)

	shellResponse, err := shellInstance.Suggest(userInput, "")
	if err != nil {
		fmt.Printf("%s: %s\n", ErrorShellAI, err)
		return
	}

	for shellResponse.CommandAction == shell.CommandActionRevise {
		shellResponse, err = shellInstance.Suggest("", shellResponse.Command)
		if err != nil {
			fmt.Printf("%s: %s\n", ErrorShellAI, err)
			return
		}
	}

	if shellResponse.CommandAction == shell.CommandActionExecute {
		err = shellInstance.Execute(shellResponse.Command)
		if err != nil {
			fmt.Printf("%s: %s\n", ErrorShellAI, err)
			return
		}
	}

}
