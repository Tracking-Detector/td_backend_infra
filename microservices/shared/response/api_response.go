package response

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Self    string      `json:"self,omitempty"`
	Next    string      `json:"next,omitempty"`
	Pages   int         `json:"pages,omitempty"`
}

// NewSuccessResponse creates a new success API response.
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse creates a new error API response.
func NewErrorResponse(message string) *APIResponse {
	return &APIResponse{
		Success: false,
		Message: message,
	}
}

func NewPagedSuccessResponse(data interface{}, self, next string, pages int) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
		Self:    self,
		Next:    next,
		Pages:   pages,
	}
}
