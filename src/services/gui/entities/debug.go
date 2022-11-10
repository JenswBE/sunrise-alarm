package entities

type DebugTemplate struct{ BaseData }

func (t DebugTemplate) GetTemplateName() string {
	return "debug"
}
