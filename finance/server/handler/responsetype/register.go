package responsetype

type RegisterSuccess struct {
	Token string `json:"token"`
}

type RegisterFail struct {
	Message string `json:"message"`
}
