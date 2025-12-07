package dto

import "time"

// HealthResponse represents the health check response data
// @Description Health check response data
type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Service   string `json:"service" example:"nutrisnap-api"`
	Version   string `json:"version" example:"1.0.0"`
	Timestamp string `json:"timestamp" example:"2024-12-07T10:00:00Z"`
}

func NewHealthResponse() HealthResponse {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return HealthResponse{
		Status:    "healthy",
		Service:   "nutrisnap-api",
		Version:   "1.0.0",
		Timestamp: time.Now().In(loc).Format(time.RFC3339),
	}
}
