package completion

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	generatePromptTemplate = `I will give you a prompt to create a single line bash command that one can enter in a terminal and run, based on what is asked in the prompt.

    {{ .details }}

    {{ .explain }}

    The prompt is: {{ .prompt }}`

	regeneratePromptTemplate = `Please update the following bash script based on what is asked in the following prompt.
    
	The script: {{ .command }}    
	The prompt: {{ .prompt }}

    {{ .details }}
	
	{{ .explain }}
	`

	promptDeteails = `Please only reply with the single line bash command surrounded by squared brackets. It should be able to be directly run in a bash terminal. Do not include any other text.`

	promptExplain = `Then please describe the bash script in plain english, step by step, what exactly it does.
  Please describe succintly, use as few words as possible, do not be verbose. 
  If there are multiple steps, please display them as a list.
`
)

var openAIResponseRegex = regexp.MustCompile("\\[(.*?)\\]\\s*(.*)")

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

func (c *Completion) Suggest(input, previousStep string) (*CompletionResponse, error) {
	if input == "" {
		return nil, fmt.Errorf("input is empty")
	}

	var prompt string
	if previousStep == "" {
		prompt = buildGenerationPrompt(input)
	} else {
		prompt = buildRenerationPrompt(input, previousStep)
	}

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
		fmt.Fprintf(os.Stderr, "🤖 OOPS!\n%s\n\n", content)
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

func buildGenerationPrompt(prompt string) string {
	var output bytes.Buffer

	templ := template.Must(template.New("prompt").Parse(generatePromptTemplate))
	err := templ.Execute(&output, map[string]interface{}{
		"details": promptDeteails,
		"explain": promptExplain,
		"prompt":  prompt,
	})
	if err != nil {
		panic(err)
	}

	return removeInitialSpaces(output.String())
}

func buildRenerationPrompt(prompt, command string) string {
	var output bytes.Buffer

	templ := template.Must(template.New("prompt").Parse(regeneratePromptTemplate))
	err := templ.Execute(&output, map[string]interface{}{
		"details": promptDeteails,
		"explain": promptExplain,
		"prompt":  prompt,
		"command": command,
	})
	if err != nil {
		panic(err)
	}

	return removeInitialSpaces(output.String())
}

func removeInitialSpaces(input string) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimLeft(line, " ")
		lines[i] = strings.TrimLeft(lines[i], "\t")
	}
	return strings.Join(lines, "\n")
}
