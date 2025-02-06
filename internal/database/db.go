package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/travismccollum/whatidid/pkg/types"
)

var db *sql.DB

func getDBPath(dbName string) string {
	// Check if running from installed location
	if os.Getenv("GOBIN") != "" || os.Getenv("GOPATH") != "" {
		dbDir := filepath.Join(os.Getenv("HOME"), ".local/share/whatidid")
		os.MkdirAll(dbDir, 0755)
		return filepath.Join(dbDir, dbName)
	}
	// Development mode - use local directory
	return dbName
}

func InitDB(dbName string) {
	dbPath := getDBPath(dbName)
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Create events table
	createEventsTable := `CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp TEXT NOT NULL,
		source TEXT NOT NULL,
		event_type TEXT NOT NULL,
		metadata TEXT,
		UNIQUE(timestamp, source, metadata)
	);`

	// Create plugin_scans table
	createScansTable := `CREATE TABLE IF NOT EXISTS plugin_scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		plugin_name TEXT NOT NULL,
		start_time TEXT NOT NULL,
		end_time TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(plugin_name, start_time, end_time)
	);`

	for _, query := range []string{createEventsTable, createScansTable} {
		if _, err := db.Exec(query); err != nil {
			log.Fatal("Failed to create table:", err)
		}
	}
}

func InsertEvent(timestamp, source, eventType, metadata string) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO events 
		(timestamp, source, event_type, metadata) 
		VALUES (?, ?, ?, ?)`,
		timestamp, source, eventType, metadata)
	return err
}

// Record a completed scan range for a plugin
func RecordPluginScan(pluginName, startTime, endTime string) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO plugin_scans 
		(plugin_name, start_time, end_time) 
		VALUES (?, ?, ?)`,
		pluginName, startTime, endTime)
	return err
}

// Get all previously scanned ranges for a plugin
func GetPluginScanRanges(pluginName string) ([]types.TimeRange, error) {
	rows, err := db.Query(`
		SELECT start_time, end_time 
		FROM plugin_scans 
		WHERE plugin_name = ? 
		ORDER BY start_time`,
		pluginName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ranges []types.TimeRange
	for rows.Next() {
		var r types.TimeRange
		if err := rows.Scan(&r.Start, &r.End); err != nil {
			return nil, err
		}
		ranges = append(ranges, r)
	}
	return ranges, nil
}

func CloseDB() {
	db.Close()
}
