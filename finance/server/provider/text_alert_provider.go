package provider

type TextAlertProvider interface {
	Alert(text string) error
}
