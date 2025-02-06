package chrome

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"log"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/travismccollum/whatidid/pkg/types"
)

type Chrome struct{
	config map[string]interface{}
}

func (c *Chrome) Name() string {
	return "chrome"
}

func (c *Chrome) Initialize(config map[string]interface{}) error {
	c.config = config
	return nil
}

// Update method signature to match interface
func (c *Chrome) FetchEvents(start, stop string, skipRanges []types.TimeRange) ([]types.Event, error) {
	// Parse input dates to time.Time
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format (use YYYY-MM-DD): %v", err)
	}

	stopTime, err := time.Parse("2006-01-02", stop)
	if err != nil {
		return nil, fmt.Errorf("invalid stop date format (use YYYY-MM-DD): %v", err)
	}

	// Add one day to stop time to include the entire day
	stopTime = stopTime.AddDate(0, 0, 1)

	// Convert to UTC format strings that SQLite can understand
	startStr := startTime.UTC().Format("2006-01-02 15:04:05")
	stopStr := stopTime.UTC().Format("2006-01-02 15:04:05")

	log.Printf("Chrome plugin: fetching history between %s and %s", startStr, stopStr)

	log.Printf("Fetching Chrome history with args: start=%s, stop=%s, skipRanges=%v", start, stop, skipRanges)
	var events []types.Event
	paths := viper.GetStringSlice("plugins.chrome.history_paths")
	stopRange := viper.GetString("plugins.chrome.stop_range")
	
	if start == "" {
		start = "2000-01-01T00:00:00Z" // Default to old date if no last fetch
	}

	log.Printf("Chrome plugin: checking %d paths, start fetch: %s, stop range: %s", len(paths), start, stopRange)

	for _, historyPath := range paths {
		// Expand home directory if path starts with ~
		if strings.HasPrefix(historyPath, "~/") {
			home, _ := os.UserHomeDir()
			historyPath = filepath.Join(home, historyPath[2:])
		}

		log.Printf("Checking Chrome history at: %s", historyPath)

		// Skip if history file doesn't exist
		if _, err := os.Stat(historyPath); os.IsNotExist(err) {
			log.Printf("History file not found: %s", historyPath)
			continue
		}

		// Copy Chrome history to temporary file (Chrome locks the original)
		tmpFile := historyPath + ".tmp"
		if err := copyFile(historyPath, tmpFile); err != nil {
			log.Printf("Error copying history file: %v", err)
			continue
		}
		defer os.Remove(tmpFile)

		// Open the temporary history file
		db, err := sql.Open("sqlite3", tmpFile)
		if err != nil {
			log.Printf("Error opening database: %v", err)
			continue
		}
		defer db.Close()

		 // Build SQL to exclude previously scanned ranges
		excludeRanges := ""
		var excludeArgs []interface{}
		
		for _, r := range skipRanges {
			if excludeRanges != "" {
				excludeRanges += " AND "
			}
			excludeRanges += `NOT (
				visits.visit_time/1000000-11644473600 >= strftime('%s', ?) 
				AND visits.visit_time/1000000-11644473600 <= strftime('%s', ?)
			)`
			excludeArgs = append(excludeArgs, r.Start, r.End)
		}

		// Query the history with proper timestamp handling and filtering
		query := `
			SELECT 
				strftime('%Y-%m-%dT%H:%M:%SZ', visits.visit_time/1000000-11644473600, 'unixepoch') as timestamp,
				urls.url,
				urls.title
			FROM visits 
			JOIN urls ON urls.id = visits.url
			WHERE visits.visit_time/1000000-11644473600 > strftime('%s', ?)
			AND (? = '' OR visits.visit_time/1000000-11644473600 <= strftime('%s', ?))`

		if excludeRanges != "" {
			query += " AND " + excludeRanges
		}

		query += " ORDER BY visits.visit_time DESC"

		args := []interface{}{startStr, stopStr, stopStr}
		args = append(args, excludeArgs...)
		
		log.Printf("Running query with lastFetch=%s stopRange=%s", start, stopRange)
		
		rows, err := db.Query(query, args...)
		if err != nil {
			log.Printf("Error querying database: %v", err)
			continue
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var timestamp, url, title string
			if err := rows.Scan(&timestamp, &url, &title); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}

			events = append(events, types.Event{
				Timestamp: timestamp,
				Source:    "chrome",
				EventType: "browser-history",
				Metadata:  title + " | " + url,
			})
			count++
		}

		if err = rows.Err(); err != nil {
			log.Printf("Error iterating rows: %v", err)
		}

		log.Printf("Found %d events in this history file", count)
	}

	log.Printf("Chrome plugin complete, found total of %d events", len(events))
	return events, nil
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if (err != nil) {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}
