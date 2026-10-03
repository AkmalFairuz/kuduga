package requesttype

type OTPInfo struct {
	ID    int64  `form:"id" validate:"required"`
	Token string `form:"token" validate:"required"`
	Code  string `form:"code" validate:"required,len=6"`
}

type ResendOTPRequest struct {
	ID    int64  `form:"id" validate:"required"`
	Token string `form:"token" validate:"required"`
}
