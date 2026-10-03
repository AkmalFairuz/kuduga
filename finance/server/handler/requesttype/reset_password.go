package requesttype

type ResetPasswordOTPRequest struct {
	Email string `form:"email" validate:"required,email,lowercase"`
}

type ResetPasswordVerifyOTPRequest struct {
	Email string  `form:"email" validate:"required,email,lowercase"`
	OTP   OTPInfo `form:"otp"`
}

type ResetPasswordRequest struct {
	ResetPasswordID int64  `form:"resetPasswordId" validate:"required"`
	Token           string `form:"token" validate:"required"`
	NewPassword     string `form:"newPassword" validate:"required,min=8,max=64"`
}
