package main

import (
	"context"
	"fmt"
	"github.com/chenjinya/data-driven-agents/agents/agent"
	"github.com/chenjinya/data-driven-agents/agents/base"
	"github.com/chenjinya/data-driven-agents/agents/pipline"
)

func main() {
	ctx := context.Background()
	agentA := agent.NewSimpleAgent(agent.SimpleAgentConfig{
		Name:            "give a number",
		InputValidator:  nil,
		PromptTemplates: []string{"请给出一个数字"},
	})

	agentB := agent.NewSimpleAgent(agent.SimpleAgentConfig{
		Name:            "it's odd or even ?",
		PromptTemplates: []string{"请判断数字 {{predict}} 是奇数还是偶数"},
	})
	agentB.SetInputValidator(func() error {
		if agentB.Input()["predict"] == nil {
			return fmt.Errorf("input is empty")
		}
		return nil
	})

	agentC := agent.NewSimpleAgent(agent.SimpleAgentConfig{
		Name:            "judgement completion is correct",
		PromptTemplates: []string{"'{{predict}}' 这个结论正确吗？"},
	})
	agentC.SetInputValidator(func() error {
		if agentC.Input()["predict"] == nil {
			return fmt.Errorf("input is empty")
		}
		return nil
	})

	agents := map[string]base.Agent{
		"entry":           agentA,
		"number is odd ?": agentB,
		"judgement":       agentC,
	}
	roadmap := pipline.RoadPath{
		ID: "entry",
		Next: []pipline.RoadPath{
			{
				ID: "number is odd ?",
				Next: []pipline.RoadPath{
					{
						ID: "judgement",
					},
				},
			},
		},
	}
	var pip base.Pipline
	pip = pipline.NewPipline(ctx, pipline.PiplineConfig{
		Name:    "guess number is odd or even",
		Agents:  agents,
		Roadmap: roadmap,
	})
	err := pip.Run()
	if err != nil {
		panic(err)
	}
}
