package requesttype

type GetTransactionsRequest struct {
	From  int64 `query:"from" validate:"required"`
	To    int64 `query:"to" validate:"required"`
	Limit int   `query:"limit" validate:"required,min=1,max=300"`
}
