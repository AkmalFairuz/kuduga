package requesttype

type CreatePurchaseRequest struct {
	ProductID      int64                             `form:"productId" validate:"required"`
	Destination_   []PurchaseDestinationFieldRequest `form:"destination"`
	Pin            string                            `form:"pin" validate:"max=6"`
	MaxPrice       int64                             `form:"maxPrice"`
	AllowDuplicate bool                              `form:"allowDuplicate"`
}

func (req CreatePurchaseRequest) Destination() map[string]string {
	return toMappedDestination(req.Destination_)
}

func toMappedDestination(destination []PurchaseDestinationFieldRequest) map[string]string {
	ret := map[string]string{}
	for _, f := range destination {
		ret[f.Key] = f.Value
	}
	return ret
}

type BillRequest struct {
	Pay    bool   `form:"pay"`
	BillID int64  `form:"id"`
	Pin    string `form:"pin" validate:"max=6"`

	ProductID    int64                             `form:"productId" validate:"required"`
	Destination_ []PurchaseDestinationFieldRequest `form:"destination"`
}

func (req BillRequest) Destination() map[string]string {
	return toMappedDestination(req.Destination_)
}

type PurchaseDestinationFieldRequest struct {
	Key   string `form:"key" validate:"required"`
	Value string `form:"value"`
}

type PurchaseHistoryRequest struct {
	BeforeID int64 `query:"beforeId"`
	Status   *int  `query:"status"`
	Limit    int64 `query:"limit" validate:"min=0,max=200"`
}

type UpdatePurchaseUserSellPriceRequest struct {
	ID            int64 `form:"id" validate:"required"`
	UserSellPrice int64 `form:"userSellPrice" validate:"required,min=1"`
}

type SearchPurchaseRequest struct {
	Query string `query:"query" validate:"required,min=3,max=50"`
}
