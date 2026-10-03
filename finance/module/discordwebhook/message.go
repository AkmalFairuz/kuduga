package discordwebhook

type Message struct {
	Content string `json:"content"`

	Username  string  `json:"username,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	Tts       bool    `json:"tts,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	Color       int    `json:"color,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
	Footer      `json:"footer,omitempty"`
	Image       Image `json:"image,omitempty"`
	Thumbnail   Image `json:"thumbnail,omitempty"`
	Author      `json:"author"`
	Fields      []Field `json:"fields,omitempty"`
}

type Footer struct {
	Text    string `json:"text,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}

type Image struct {
	Url string `json:"url,omitempty"`
}

type Author struct {
	Name    string `json:"name,omitempty"`
	Url     string `json:"url,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}

type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}
