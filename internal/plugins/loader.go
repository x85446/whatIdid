package plugins

import (
	"log"

	"github.com/spf13/viper"
	"github.com/travismccollum/whatidid/pkg/types"
	"github.com/travismccollum/whatidid/internal/plugins/chrome"
)

func LoadEnabledPlugins() []types.Plugin {
	viper.SetConfigName("wid")
	// Look for config in project directory first, then fall back to user's home
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.config/whatidid")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: Could not read config file, using defaults: %v", err)
	}

	var activePlugins []types.Plugin

	// Always load the Dummy plugin for testing
	activePlugins = append(activePlugins, &Dummy{})

	if viper.GetBool("plugins.chrome.enabled") {
		activePlugins = append(activePlugins, &chrome.Chrome{})
	}

	// Commented out until plugins are implemented
	/*
	if viper.GetBool("plugins.chrome.enabled") {
		activePlugins = append(activePlugins, &Chrome{})
	}
	if viper.GetBool("plugins.git.enabled") {
		activePlugins = append(activePlugins, &Git{})
	}
	if viper.GetBool("plugins.outlook.enabled") {
		activePlugins = append(activePlugins, &Outlook{})
	}
	*/

	return activePlugins
}
