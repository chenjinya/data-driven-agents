package base

type JSON map[string]interface{}

func (j JSON) Validate() error {
	return nil
}

func (j JSON) Predict() string {
	if j["predict"] == nil {
		return ""
	}
	return j["predict"].(string)
}

type Agent interface {
	Name() string
	Prompt(...string) string
	Input() JSON
	SetInput(JSON) error
	SetInputValidator(func() error)
	Output() JSON
	Call() error
}

type AgentRoad struct {
	Agent Agent
	Road  map[string]*AgentRoad
}

type Pipline interface {
	Name() string
	Execute() error
}
