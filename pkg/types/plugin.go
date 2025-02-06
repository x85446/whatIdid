package types

type TimeRange struct {
	Start string
	End   string
}

type Event struct {
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
	EventType string `json:"event_type"`
	Metadata  string `json:"metadata"`
}

type Plugin interface {
	Name() string
	Initialize(config map[string]interface{}) error
	// Updated to include ranges to skip
	FetchEvents(start, stop string, skipRanges []TimeRange) ([]Event, error)
}
