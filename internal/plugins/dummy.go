package plugins

import (
	"github.com/travismccollum/whatidid/pkg/types"
)

type Dummy struct{}

func (d *Dummy) Name() string {
	return "dummy"
}

func (d *Dummy) Initialize(config map[string]interface{}) error {
	return nil
}

func (d *Dummy) FetchEvents(start, stop string, skipRanges []types.TimeRange) ([]types.Event, error) {
	// Just return an empty slice for the dummy plugin
	return []types.Event{}, nil
}
