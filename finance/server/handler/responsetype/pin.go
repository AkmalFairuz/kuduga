package responsetype

type IsPinEnabledResponse struct {
	Enabled bool `json:"enabled"`
}

type ValidatePinResponse struct {
	Valid bool `json:"valid"`
}
