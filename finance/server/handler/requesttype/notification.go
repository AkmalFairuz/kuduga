package requesttype

type UpdateFcmTokenRequest struct {
	FcmToken string `form:"fcmToken" validate:"required"`
}
