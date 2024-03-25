package agent

import (
	"context"
	"fmt"
	"github.com/chenjinya/data-driven-agents/agents/base"
	"github.com/hoisie/mustache"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
)

type SimpleAgent struct {
	name           string
	system         string
	input          base.JSON
	inputValidator func() error
	output         base.JSON
	promptTemplate string
}

type SimpleAgentConfig struct {
	Name            string
	System          string
	InputValidator  func() error
	PromptTemplates []string
}

func NewSimpleAgent(cfg SimpleAgentConfig) base.Agent {
	return &SimpleAgent{
		name:           cfg.Name,
		system:         cfg.System,
		inputValidator: cfg.InputValidator,
		promptTemplate: strings.Join(cfg.PromptTemplates, "\n"),
	}
}
func (a *SimpleAgent) Name() string {
	if a.name == "" {
		return "SimpleAgent"
	}
	return a.name
}
func (a *SimpleAgent) Prompt(promptTemplates ...string) string {
	if len(promptTemplates) > 0 {
		a.promptTemplate = strings.Join(promptTemplates, "\n")
	}
	return a.promptTemplate
}

func (a *SimpleAgent) System() string {
	if a.system == "" {
		return "You are powerful assistant"
	}
	return a.system
}

func (a *SimpleAgent) SetInput(input base.JSON) error {
	a.input = input
	if err := a.inputValidator(); err != nil {
		return err
	}
	return nil
}
func (a *SimpleAgent) Input() base.JSON {
	return a.input
}

func (a *SimpleAgent) SetInputValidator(v func() error) {
	a.inputValidator = v
}

func (a *SimpleAgent) Output() base.JSON {
	return a.output
}
func (a *SimpleAgent) Call(ctx context.Context) (err error) {
	cfg := openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))
	cfg.BaseURL = os.Getenv("OPENAI_BASE_URL")
	client := openai.NewClientWithConfig(cfg)
	mustmp, _ := mustache.ParseString(a.Prompt())
	prompt := mustmp.Render(a.Input())
	fmt.Println("🤖 Predict prompt: ", prompt)
	result, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: a.System(),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("🤖 ChatCompletion error: %v\n", err)
		return err
	}

	a.output = base.JSON{
		"predict": result.Choices[0].Message.Content,
	}
	return nil
}
