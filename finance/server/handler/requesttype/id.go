package requesttype

type IDRequest struct {
	ID int64 `json:"id" form:"id" query:"id" validate:"required"`
}
