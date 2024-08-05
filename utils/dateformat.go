package utils

import (
	"fmt"
	"time"
)

func Datefromating(unformattedDate string) (string, error) {
	fmt.Println("Received date:", unformattedDate)

	// Handle different formats
	var formats = []string{
		"2006-01-02",               // YYYY-MM-DD
		"2006-01-02T15:04:05Z",     // ISO 8601 with time
		"2006-01-02T15:04:05.000Z", // ISO 8601 with milliseconds
		"2006-01-02T15:04:05.000",  // ISO 8601 with milliseconds
	}

	for _, format := range formats {
		dt, err := time.Parse(format, unformattedDate)
		if err == nil {
			return dt.Format("2006-01-02"), nil // Format for PostgreSQL date type
		}
	}

	return "", fmt.Errorf("error parsing date: %s", unformattedDate)
}

// func Datefromating(unformattedDate string) string {
// 	fmt.Println("Received date:", unformattedDate)

// 	// Handle different formats
// 	// Assuming input might be in YYYY-MM-DD or ISO 8601 format
// 	var formats = []string{
// 		"2006-01-02",           // YYYY-MM-DD
// 		"2006-01-02T15:04:05Z", // ISO 8601 with time
// 	}

// 	for _, format := range formats {
// 		dt, err := time.Parse(format, unformattedDate)
// 		if err == nil {
// 			return dt.Format("2006-01-02") // Format for PostgreSQL date type
// 		}
// 	}

// 	fmt.Println("Error parsing date:", unformattedDate)
// 	return "Error parsing date"
// }

// func Datefromating(unformattedDate string) (string, error) {
// 	print(unformattedDate)
// 	if unformattedDate == "" {
// 		return "", nil // Return empty string if the date is not provided
// 	}
// 	// Parse ISO 8601 formatted string to time.Time object
// 	dt, err := time.Parse("2006-01-02T15:04:05.000", unformattedDate)
// 	if err != nil {
// 		// If it fails, try parsing with only date part
// 		dt, err = time.Parse("2006-01-02", unformattedDate)
// 		if err != nil {
// 			fmt.Println("Error parsing date:", err)
// 			return "", fmt.Errorf("error parsing date: %v", err)
// 		}
// 	}
// 	// Format to YYYY-MM-DD format for PostgreSQL date type
// 	formattedDate := dt.Format("2006-01-02")
// 	return formattedDate, nil
// }
