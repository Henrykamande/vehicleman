package utils

import (
	"database/sql"
	"errors"
)

// getOwnerID fetches the owner_id based on the provided email.
func GetOwnerID(db *sql.DB, email string) (int, error) {
	var ownerID int
	err := db.QueryRow("SELECT user_id FROM users WHERE email = $1", email).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}
		return 0, err
	}
	return ownerID, nil
}
