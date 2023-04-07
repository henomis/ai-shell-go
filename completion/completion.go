package completion

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	promptTemplate = `I will give you a prompt to create a single line bash command that one can enter in a terminal and run, based on what is asked in the prompt.

    {{ .details }}

    {{ .explain }}

    The prompt is: {{ .prompt }}`

	details = `Please only reply with the single line bash command surrounded by 3 backticks. It should be able to be directly run in a bash terminal. Do not include any other text.`

	explain = `Then please describe the bash script in plain english, step by step, what exactly it does.
  Please describe succintly, use as few words as possible, do not be verbose. 
  If there are multiple steps, please display them as a list.
`
)

var openAIResponseRegex = regexp.MustCompile("```(.*?)```\\s*(.*)")

type Completion struct {
	openAIClient *openai.Client
}

type CompletionResponse struct {
	Command string
	Explain string
}

func New(openAIClient *openai.Client) *Completion {
	return &Completion{
		openAIClient: openAIClient,
	}
}

func (c *Completion) Suggest(input string) (*CompletionResponse, error) {
	if input == "" {
		return nil, fmt.Errorf("input is empty")
	}

	prompt := buildPrompt(details, explain, input)

	response, err := c.openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	content := response.Choices[0].Message.Content
	matches := openAIResponseRegex.FindStringSubmatch(content)

	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("no command found")
	}

	command := matches[1]
	explanation := strings.TrimSpace(matches[2])

	return &CompletionResponse{
		Command: command,
		Explain: explanation,
	}, nil

}

// ---------------
// support methods
// ---------------

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
