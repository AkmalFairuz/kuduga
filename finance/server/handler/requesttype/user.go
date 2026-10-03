package requesttype

type UpdateUserEmailRequest struct {
	OTP   OTPInfo `form:"otp" validate:"required"`
	Email string  `form:"email" validate:"required,email"`
}

type UpdateUserDisplayNameRequest struct {
	DisplayName string `form:"displayName" validate:"required"`
}

type UpdateUserEmailOTPRequest struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

type UpdateUserPasswordRequest struct {
	CurrentPassword string `form:"currentPassword" validate:"required"`
	NewPassword     string `form:"newPassword" validate:"required,min=8,max=64"`
}

type KycRequest struct {
	DocumentID string `form:"documentId" validate:"required"`
	FullName   string `form:"fullName" validate:"required,max=255"`
}

type UseReferralCodeRequest struct {
	ReferralCode string `form:"referralCode" validate:"required"`
}

type UpdateUserDeviceUniqueIDRequest struct {
	DeviceUniqueID string `form:"deviceUniqueId" validate:"required,max=255"`
}
