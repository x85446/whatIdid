package database

import (
	"database/sql"
	"log"
)

func GetLastFetchTime() string {
	var lastFetchTime string
	err := db.QueryRow("SELECT MAX(timestamp) FROM events").Scan(&lastFetchTime)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting last fetch timestamp: %v\n", err)
		return ""
	}
	return lastFetchTime
}
