package responsetype

type NotificationResponse struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Type        int               `json:"type"`
	HasRead     bool              `json:"hasRead"`
	RefID       string            `json:"refId"`
	Data        map[string]string `json:"data"`
	CreatedAt   int64             `json:"createdAt"`
}
