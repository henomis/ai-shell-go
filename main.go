package main

import (
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"github.com/henomis/ai-shell-go/completion"
	"github.com/henomis/ai-shell-go/shell"
)

var (
	ErrorShellAI = fmt.Errorf("something went wrong")
)

func main() {

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		fmt.Println("OPEN_AI_KEY is not set.")
		fmt.Println("Please set the OPENAI_API_KEY environment variable to your OpenAI API key.")
		return
	}

	userInput := strings.Join(os.Args[1:], " ")

	client := openai.NewClient(openAIKey)
	completion := completion.New(client)
	s := shell.New(completion)

	shellResponse, err := s.Suggest(userInput)
	if err != nil {
		fmt.Printf("%s: %s\n", ErrorShellAI, err)
		return
	}

	for shellResponse.CommandAction == shell.CommandActionRevise {
		shellResponse, err = s.Retry(shellResponse.Command)
		if err != nil {
			fmt.Printf("%s: %s\n", ErrorShellAI, err)
			return
		}
	}

	if shellResponse.CommandAction == shell.CommandActionExecute {
		err = s.Execute(shellResponse.Command)
		if err != nil {
			fmt.Printf("%s: %s\n", ErrorShellAI, err)
			return
		}
	}

}
