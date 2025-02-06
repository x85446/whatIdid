package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/travismccollum/whatidid/internal/database"
	"github.com/travismccollum/whatidid/internal/plugins"
	"github.com/travismccollum/whatidid/pkg/types"  // Add this import
)

var (
	startTime string
	stopTime  string
)

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVar(&startTime, "start", "", "Start time (YYYY-MM-DD)")
	fetchCmd.Flags().StringVar(&stopTime, "stop", "", "Stop time (YYYY-MM-DD)")
}

var fetchCmd = &cobra.Command{
	Use:   "fetch [plugin]",
	Short: "Fetch events from enabled plugins and store them in the database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Fetching latest events...")

		// Initialize the database
		database.InitDB("whatidid.db")
		defer database.CloseDB()

		 // Set default time range if not specified
		if startTime == "" {
			lastFetchTime := database.GetLastFetchTime()
			if lastFetchTime == "" {
				startTime = time.Now().AddDate(0, -1, 0).Format("2006-01-02") // default to 1 month ago
			} else {
				startTime = lastFetchTime
			}
		}
		if stopTime == "" {
			stopTime = time.Now().Format("2006-01-02")
		}

		// Load specific plugin if specified
		var targetPlugins []types.Plugin
		if len(args) > 0 {
			pluginName := args[0]
			allPlugins := plugins.LoadEnabledPlugins()
			for _, p := range allPlugins {
				if p.Name() == pluginName {
					targetPlugins = append(targetPlugins, p)
					break
				}
			}
			if len(targetPlugins) == 0 {
				log.Fatalf("Plugin %s not found", pluginName)
			}
		} else {
			targetPlugins = plugins.LoadEnabledPlugins()
		}

		for _, plugin := range targetPlugins {
			fmt.Printf("Fetching from plugin: %s\n", plugin.Name())
			
			 // Get previously scanned ranges for this plugin
			skipRanges, err := database.GetPluginScanRanges(plugin.Name())
			if err != nil {
				log.Printf("Error getting scan ranges for %s: %v\n", plugin.Name(), err)
				continue
			}

			// Initialize plugin with empty config for now
			if err := plugin.Initialize(map[string]interface{}{}); err != nil {
				log.Printf("Error initializing plugin %s: %v\n", plugin.Name(), err)
				continue
			}

			events, err := plugin.FetchEvents(startTime, stopTime, skipRanges)
			if err != nil {
				log.Printf("Error fetching from plugin %s: %v\n", plugin.Name(), err)
				continue
			}

			fmt.Printf("Found %d events from %s\n", len(events), plugin.Name())

			// Insert events and record the scan range
			if len(events) > 0 {
				var minTime, maxTime string
				minTime = events[0].Timestamp
				maxTime = events[0].Timestamp

				for _, event := range events {
					if event.Timestamp < minTime {
						minTime = event.Timestamp
					}
					if event.Timestamp > maxTime {
						maxTime = event.Timestamp
					}

					if err := database.InsertEvent(event.Timestamp, event.Source, 
						event.EventType, event.Metadata); err != nil {
						log.Printf("Failed to insert event: %v\n", err)
					}
				}

				// Record the scanned range
				if err := database.RecordPluginScan(plugin.Name(), minTime, maxTime); err != nil {
					log.Printf("Failed to record scan range: %v\n", err)
				}
			}
		}

		fmt.Println("Fetch complete.")
	},
}
