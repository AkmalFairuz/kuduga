package requesttype

type CheckerRequest struct {
	ID    string `query:"id" validate:"required"`
	Input []struct {
		Key   string `query:"key" validate:"required"`
		Value string `query:"value"`
	} `query:"input" validate:"required"`
}
