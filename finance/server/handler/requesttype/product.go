package requesttype

type GetProductCategoriesRequest struct {
	IsNested bool  `query:"isNested"`
	ParentID int64 `query:"parentId"`
}
