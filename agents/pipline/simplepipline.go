package pipline

import (
	"context"
	"fmt"
	"github.com/chenjinya/data-driven-agents/agents/base"
)

type PiplineConfig struct {
	Name    string
	Agents  map[string]base.Agent
	Roadmap RoadPath
}

type RoadPath struct {
	ID        string
	Condition string
	Next      []RoadPath
}

type SimplePipline struct {
	ctx           context.Context
	name          string
	roadmap       RoadPath
	agents        map[string]base.Agent
	nextroads     []RoadPath
	prevagent     base.Agent
	BeforeStart   func(data any) error
	BeforeExecute func(data any) error
	AfterExecuted func(data any) error
	AfterFinished func(data any) error
}

func (p *SimplePipline) Name() string {
	return "simple pipline"
}

func NewPipline(ctx context.Context, cfg PiplineConfig) base.Pipline {
	return &SimplePipline{
		ctx:     ctx,
		name:    cfg.Name,
		roadmap: cfg.Roadmap,
		agents:  cfg.Agents,
		BeforeStart: func(data any) error {
			fmt.Println("🚀 Pipline:", data.(base.Pipline).Name())
			return nil
		},
		BeforeExecute: func(data any) error {
			fmt.Println("🚀 Agent:", data.(RoadPath).ID)
			return nil
		},
		AfterExecuted: func(result any) error {
			fmt.Println("🚀 Result:", result.(ExecutedResult).Result.Predict())
			return nil
		},
		AfterFinished: func(data any) error {
			fmt.Println("🚀 Finished")
			return nil
		},
	}
}
func (p *SimplePipline) event(fn func(d any) error, data any) error {
	if fn != nil {
		return fn(data)
	}
	return nil
}
func (p *SimplePipline) Run() error {
	if err := p.event(p.BeforeStart, p); err != nil {
		return err
	}
	if err := p.Execute(); err != nil {
		return err
	}
	if err := p.event(p.AfterFinished, nil); err != nil {
		return err
	}
	return nil
}

type ExecutedResult struct {
	Err    error
	Result base.JSON
}

func (p *SimplePipline) Execute() error {
	if p.nextroads == nil {
		p.nextroads = []RoadPath{p.roadmap}
	}
	for _, item := range p.nextroads {
		if err := p.event(p.BeforeExecute, item); err != nil {
			return err
		}
		_agent := p.agents[item.ID]
		if p.prevagent != nil {
			_ = _agent.SetInput(p.prevagent.Output())
		}
		err := _agent.Call(p.ctx)
		result := ExecutedResult{
			Err: err, Result: _agent.Output(),
		}
		if err := p.event(p.AfterExecuted, result); err != nil {
			return err
		}
		if err != nil {
			return err
		}
		if item.Next != nil {
			p.prevagent = _agent
			p.nextroads = item.Next
			err = p.Execute()
			if err != nil {
				return err
			}
		}

	}
	return nil
}
