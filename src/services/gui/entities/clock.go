package entities

type ClockTemplate struct{ BaseData }

func (t ClockTemplate) GetTemplateName() string {
	return "clock"
}
