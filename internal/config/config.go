package config

import (
	"encoding/json"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/wo0lien/cosmoBot/internal/logging"
)

var Config *ConfigStruct

// Initialising the config
func init() {
	cfg, err := LoadConfig()
	if err != nil {
		logging.Error.Printf("Fatal error loading config : %s", err)
		return
	}
	Config = cfg
}

type EventType string

const (
	EventTypeVisit         EventType = "CosmoVisit"
	EventTypeApero         EventType = "CosmoApero"
	EventTypeBaO           EventType = "CosmoBaO"
	EvetTypeCafeDesLangues EventType = "CafeDesLangues"
	EventTypePerm          EventType = "Perm"
	EventTypeOther         EventType = "Other"
)

type ResponseMethod struct {
	DiscordId   string                `json:"discordId"`
	ChannelType discordgo.ChannelType `json:"channelType"`
}

//go:generate go run schema.gen.go

type ConfigStruct struct {
	schema                    string                       `json:"$schema",omitempty`
	ResponseMethodByEventType map[EventType]ResponseMethod `json:"responseMethodByEventType"`
}

// load config from config.json file in the root directory
func LoadConfig() (*ConfigStruct, error) {
	var config ConfigStruct

	// check if config.json exists
	_, err := os.Stat("config.json")

	if err != nil {
		return nil, err
	}

	// create config.json with default values
	f, err := os.ReadFile("config.json")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(f, &config)

	if err != nil {
		return nil, err
	}

	logging.Info.Println("Config loaded")

	return &config, nil
}

// Save config to config.json file in the root directory
func SaveConfig(config *ConfigStruct) error {
	configBytes, err := json.MarshalIndent(config, "", "  ")

	if err != nil {
		return err
	}

	err = os.WriteFile("config.json", configBytes, 0644)

	if err != nil {
		return err
	}

	return nil
}