package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	openai "github.com/sashabaranov/go-openai"

	"github.com/henomis/ai-shell-go/completion"
	"github.com/henomis/ai-shell-go/shell"
)

const promptTemplate = `I will give you a prompt to create a single line bash command that one can enter in a terminal and run, based on what is asked in the prompt.

    {{ .details }}

    {{ .explain }}

    The prompt is: {{ .prompt }}`

const details = `Please only reply with the single line bash command surrounded by 3 backticks. It should be able to be directly run in a bash terminal. Do not include any other text.`

const explain = `Then please describe the bash script in plain english, step by step, what exactly it does.
  Please describe succintly, use as few words as possible, do not be verbose. 
  If there are multiple steps, please display them as a list.
`

func buildPrompt(details, explain, prompt string) string {
	var output bytes.Buffer

	templ := template.Must(template.New("prompt").Parse(promptTemplate))
	err := templ.Execute(&output, map[string]interface{}{
		"details": details,
		"explain": explain,
		"prompt":  prompt,
	})
	if err != nil {
		panic(err)
	}

	return output.String()
}

func main() {

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		fmt.Println("OPEN_AI_KEY is not set")
		return
	}

	client := openai.NewClient(openAIKey)
	completion := completion.New(client)
	s := shell.New(completion)

	s.Run()
}
