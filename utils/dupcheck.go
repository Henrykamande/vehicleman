package utils

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func duplicateCheck(db *sql.DB, valuetocheck string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Owners WHERE owner_contact = $1", valuetocheck).Scan(&count)
	if err != nil {
		// Handle error appropriately, e.g., log or return false
		return false
	}
	return count > 0
}
