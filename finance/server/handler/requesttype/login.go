package requesttype

type LoginRequest struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required,max=50"`
}

type LoginWithGoogleRequest struct {
	IDToken string `form:"idToken" validate:"required"`
}
