package responsetype

type ResetPasswordOTPResponse struct {
	OTP OTPInfo `json:"otp"`
}

type ResetPasswordVerifyOTPResponse struct {
	ResetPasswordID int64  `json:"resetPasswordId"`
	Token           string `json:"token"`
}
