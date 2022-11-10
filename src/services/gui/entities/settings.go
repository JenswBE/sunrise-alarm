package entities

type SettingsTemplate struct{ BaseData }

func (t SettingsTemplate) GetTemplateName() string {
	return "settings"
}
