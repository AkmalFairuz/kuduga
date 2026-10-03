package requesttype

type ValidatePinRequest struct {
	Pin string `form:"pin" validate:"required,number,len=6"`
}

type CreatePinRequest struct {
	Pin string `form:"pin" validate:"required,number,len=6"`
}

type DeletePinRequest struct {
	Pin string `form:"pin" validate:"required,number,len=6"`
}

type UpdatePinRequest struct {
	Pin    string `form:"pin" validate:"required,number,len=6"`
	NewPin string `form:"newPin" validate:"required,number,len=6"`
}
