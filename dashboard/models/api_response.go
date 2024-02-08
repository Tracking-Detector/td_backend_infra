package models

type APIResponse[T interface{}] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Self    string `json:"self,omitempty"`
	Next    string `json:"next,omitempty"`
	Pages   int    `json:"pages,omitempty"`
}
