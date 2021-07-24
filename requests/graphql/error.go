package graphql

type Error interface {
	error
	Network() bool // Is the error a network error?
	Server() bool  // Is the server error?
}

type GraphqlError struct {
	Message   string `json:"message"`
	Locations []struct {
		Line   int `json:"line"`
		Column int `json:"column"`
	} `json:"locations"`
	Extensions struct {
		Code string `json:"code"`
	} `json:"extensions"`
}
