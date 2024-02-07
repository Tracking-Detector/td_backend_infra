package models

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Self    string      `json:"self,omitempty"`
	Next    string      `json:"next,omitempty"`
	Pages   int         `json:"pages,omitempty"`
}
