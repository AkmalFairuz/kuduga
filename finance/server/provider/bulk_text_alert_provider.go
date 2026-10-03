package provider

type BulkTextAlertProvider struct {
	providers []TextAlertProvider
}

func NewBulkTextAlertSender(providers []TextAlertProvider) *BulkTextAlertProvider {
	for _, provider := range providers {
		if _, ok := provider.(*BulkTextAlertProvider); ok {
			panic("cannot create BulkTextAlertProvider with BulkTextAlertProvider in providers")
		}
	}
	return &BulkTextAlertProvider{providers: providers}
}

func (p BulkTextAlertProvider) Alert(text string) error {
	for _, provider := range p.providers {
		_ = provider.Alert(text)
	}
	return nil
}
