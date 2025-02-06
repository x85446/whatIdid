package database

import (
	"fmt"
	"time"

	"github.com/travismccollum/whatidid/pkg/types"
)

func QueryEvents(startTime, endTime string, source string) ([]types.Event, error) {
	query := "SELECT timestamp, source, event_type, metadata FROM events WHERE 1=1"
	var args []interface{}

	if startTime != "" {
		query += " AND timestamp >= ?"
		args = append(args, startTime)
	}
	if endTime != "" {
		query += " AND timestamp <= ?"
		args = append(args, endTime)
	}
	if source != "" {
		query += " AND source = ?"
		args = append(args, source)
	}

	query += " ORDER BY timestamp DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []types.Event
	for rows.Next() {
		var e types.Event
		if err := rows.Scan(&e.Timestamp, &e.Source, &e.EventType, &e.Metadata); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func ApplyRetentionPolicy(daysToKeep int) error {
	if daysToKeep <= 0 {
		return fmt.Errorf("invalid retention period: %d days", daysToKeep)
	}

	cutoff := time.Now().AddDate(0, 0, -daysToKeep).Format(time.RFC3339)
	_, err := db.Exec("DELETE FROM events WHERE timestamp < ?", cutoff)
	return err
}
