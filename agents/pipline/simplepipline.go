package pipline

import (
	"fmt"
	"github.com/chenjinya/data-driven-agents/agents/base"
)

type PiplineConfig struct {
	Name    string
	Agents  map[string]base.Agent
	Roadmap base.JSON
}

type SimplePipline struct {
	name      string
	roadmap   base.JSON
	agents    map[string]base.Agent
	nextroads base.JSON
	prevagent base.Agent
}

func (p *SimplePipline) Name() string {
	return "gus num"
}

func NewPipline(cfg PiplineConfig) base.Pipline {
	return &SimplePipline{
		name:    cfg.Name,
		roadmap: cfg.Roadmap,
		agents:  cfg.Agents,
	}
}

func (p *SimplePipline) Execute() error {
	if p.nextroads == nil {
		p.nextroads = p.roadmap
	}
	for name, next := range p.nextroads {
		_agent := p.agents[name]
		if p.prevagent != nil {
			_ = _agent.SetInput(p.prevagent.Output())
		}
		fmt.Println("> Execute agent: ", _agent.Name())
		err := _agent.Call()
		if err != nil {
			return err
		}
		result := _agent.Output()
		fmt.Println("> Result:\n", result.Predict())
		if next != nil {
			p.prevagent = _agent
			p.nextroads = next.(base.JSON)
			err = p.Execute()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
