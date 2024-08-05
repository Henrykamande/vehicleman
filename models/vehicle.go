package models

// Lorry represents a lorry with its details
type Vehicle struct {
	VehicleID          int    `json:"vehicle_id"`
	Make               string `json:"make"`
	Model              string `json:"model"`
	Year               int    `json:"year"`
	RegistrationNumber string `json:"registration_number"`
	Capacity           int    `json:"capacity"`
	AvailabilityStatus string `json:"availability_status"`
	OwnerID            int    `json:"owner_id"`
}
