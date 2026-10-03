package requesttype

type RegisterRequest struct {
	Name           string  `form:"name" validate:"required,lowercase,min=3,max=20"`
	DisplayName    string  `form:"displayName" validate:"required,min=3,max=50"`
	Email          string  `form:"email" validate:"required,email,lowercase,max=255"`
	Password       string  `form:"password" validate:"required,min=8,max=50"`
	DeviceUniqueID *string `form:"deviceUniqueId" validate:"max=255"`
	OTP            OTPInfo `form:"otp" validate:"required"`
}

type RegisterOTPRequest struct {
	Name           string  `form:"name" validate:"required,lowercase,min=3,max=20"`
	DisplayName    string  `form:"displayName" validate:"required,min=3,max=50"`
	Email          string  `form:"email" validate:"required,email,lowercase,max=255"`
	DeviceUniqueID *string `form:"deviceUniqueId" validate:"max=255"`
	Password       string  `form:"password" validate:"required,min=8,max=50"`
}

type RegisterWithGoogleRequest struct {
	IDToken        string  `form:"idToken" validate:"required"`
	DeviceUniqueID *string `form:"deviceUniqueId" validate:"max=255"`
}
