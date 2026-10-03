package responsetype

type AppLayoutInfoResponse struct {
	Announcement string `json:"announcement"`
	Warning      string `json:"warning"`
	Services     []struct {
		Title    string `json:"title"`
		Services []struct {
			Name       string            `json:"name"`
			ActionType string            `json:"actionType"`
			ActionData map[string]string `json:"actionData"`
			IconType   string            `json:"iconType"`
			IconData   string            `json:"iconData"`
		} `json:"services"`
	} `json:"services"`
}

type AppBannerInfoResponse struct {
	Image string `json:"image"`
	Url   string `json:"url"`
}

type AppVersionResponse struct {
	Version        int `json:"version"`
	MinimumVersion int `json:"minimumVersion"`
}
