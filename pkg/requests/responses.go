package requests

//easyjson:json
type (
	Response struct {
		Status int `json:"status"`
		Body   any `json:"body"`
	}
)
