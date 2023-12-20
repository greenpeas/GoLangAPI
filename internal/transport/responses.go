package transport

type SuccessResponse struct {
	Success bool `json:"success"`
}

type DeleteResponse struct {
	Success bool `json:"success"`
}

type AddTelemetryResponse struct {
	CountAdded int `json:"count_added"`
}
