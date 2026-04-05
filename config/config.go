package config

import (
	"encoding/json"
	"os"
)

const rootPath = "./config.json"

// Config represents application parameters
type Config struct {
	Label     string // Specific label to backup
	Anonymize bool   // if false we save file with date_object
	Port      int64  // Port to launch server to wait authentication code in google service

	OutputDir string `json:"output_dir"`

	MaxResults int64 `json:"max_results"`

	//Since           string `json:"since"` // TODO a voir avec un flags
	CredentialsFile string `json:"credentials_file"`
	TokenFile       string `json:"token_file"`
	Path            string
}

// Init Create the config file the first time
func Init() *Config {
	config := Config{
		OutputDir:       "./backup",
		Label:           "INBOX",
		MaxResults:      5000,
		CredentialsFile: "credentials.json",
		TokenFile:       "token.json",
		Path:            rootPath,
	}

	data, _ := json.Marshal(config)
	os.WriteFile(config.Path, data, 0600)

	return &config
}

// Load Deserialize configuration
func Load() (*Config, error) {
	config := Config{}
	file, err := os.ReadFile(rootPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &config)
	return &config, err
}

// Save Serialize configuration
func (cfg *Config) Save() error {
	data, _ := json.Marshal(cfg)
	return os.WriteFile(cfg.Path, data, 0600)
}
