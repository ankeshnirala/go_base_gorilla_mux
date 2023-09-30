package types

import "net/http"

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type StandardResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
