package responsetype

type UserKycStatusResponse struct {
	CurrentStatus bool `json:"currentStatus"`
	RequestStatus *int `json:"requestStatus"`
}

type UserReferralInfoResponse struct {
	Message                        string `json:"message"`
	YourReferralCode               string `json:"yourReferralCode"`
	TotalUserUsingYourReferralCode int64  `json:"totalUserUsingYourReferralCode"`
	ShareText                      string `json:"shareText"`
}
