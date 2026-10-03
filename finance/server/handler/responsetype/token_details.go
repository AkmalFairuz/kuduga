package responsetype

type TokenDetails struct {
	UserID          int64  `json:"userId"`
	UserName        string `json:"userName"`
	UserDisplayName string `json:"userDisplayName"`
	UserEmail       string `json:"userEmail"`
	UserRole        int    `json:"userRole"`
	HasPin          bool   `json:"hasPin"`
}
