package main

import (
	"fmt"
	"github.com/chenjinya/data-driven-agents/agents/agent"
	"github.com/chenjinya/data-driven-agents/agents/base"
	"github.com/chenjinya/data-driven-agents/agents/pipline"
)

func main() {
	agentA := agent.NewSimpleAgent(agent.SimpleAgentConfig{
		Name:            "guess a number",
		InputValidator:  nil,
		PromptTemplates: []string{"请给出一个数字"},
	})

	agentB := agent.NewSimpleAgent(agent.SimpleAgentConfig{
		Name:            "guess a number",
		PromptTemplates: []string{"请判断数字 {{predict}} 是奇数还是偶数"},
	})

	agentB.SetInputValidator(func() error {
		if agentB.Input()["predict"] == nil {
			return fmt.Errorf("input is empty")
		}
		return nil
	})
	agents := map[string]base.Agent{
		"entry":           agentA,
		"number is odd ?": agentB,
	}
	roadmap := base.JSON{
		"entry": base.JSON{
			"number is odd ?": nil,
		},
	}
	var pip base.Pipline
	pip = pipline.NewPipline(pipline.PiplineConfig{
		Agents:  agents,
		Roadmap: roadmap,
	})
	pip.Execute()
}
