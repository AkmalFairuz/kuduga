package requesttype

type AppHeader struct {
	Version int64 `reqHeader:"X-App-Version"`
}
