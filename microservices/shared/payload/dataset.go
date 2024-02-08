package payload

type CreateDatasetPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Label       string `json:"label"`
}
