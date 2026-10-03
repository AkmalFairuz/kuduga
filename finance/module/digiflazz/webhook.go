package digiflazz

import "encoding/json"

type WebhookEvent struct {
	Event string
	Data  json.RawMessage
}
