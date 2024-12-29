package types

type Err struct {
	Id    string `json:"id,omitempty"`
	Error string `json:"error"`
}

type CreateId struct {
	Id string `json:"id"`
}
