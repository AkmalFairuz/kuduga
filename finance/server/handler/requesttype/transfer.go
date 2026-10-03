package requesttype

type BalanceTransferRequest struct {
	DestinationEmail string `form:"destinationEmail" validate:"required"`
	Amount           int64  `form:"amount" validate:"required,number,min=1,max=999999999"`
	Note             string `form:"note" validate:"max=255"`
	Purpose          int    `form:"purpose" validate:"required"`
	Pin              string `form:"pin" validate:"omitempty,len=6,number"`
}

type BalanceTransferToUserCheckDestinationRequest struct {
	Email string `form:"email" validate:"required,email"`
}
