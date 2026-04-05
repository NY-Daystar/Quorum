package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"quorum/config"
	"runtime"
	"strings"
)

// GetName return filename only without extension
// ex: sample.pdf will return sample
func GetName(filePath string) string {
	filename := filepath.Base(filePath)

	var extension = filepath.Ext(filename)
	return filename[0 : len(filename)-len(extension)]
}

// GetAppDataPath returns the path to the AppData directory on Windows
// and the home directory on Linux.
func GetAppDataPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		appDataPath := os.Getenv("APPDATA")
		if appDataPath == "" {
			return "", fmt.Errorf("APPDATA environment variable is not set")
		}
		return appDataPath, nil
	} else {
		return homeDir, nil
	}
}

// GetLogsFile get logs filePath
func GetLogsFile() string {
	appDataPath, _ := GetAppDataPath()
	var logsFolder = path.Join(appDataPath, config.AppName, "logs")
	return path.Join(logsFolder, "log.json")
}

// SanitizeFilename change name to replace
func SanitizeFilename(s string) string {
	subsitute := ""
	replacer := strings.NewReplacer(
		"/", subsitute,
		"\\", subsitute,
		":", subsitute,
		"*", subsitute,
		"?", subsitute,
		"!", subsitute,
		"\"", subsitute,
		"<", subsitute,
		">", subsitute,
		"|", subsitute,
		".", subsitute,
		",", subsitute,
		"º", subsitute,
		"(", subsitute,
		")", subsitute,
		"'", " ",
	)

	result := strings.TrimSpace(replacer.Replace(s))

	maxBounds := 180

	if len(result) > maxBounds {
		result = result[0:maxBounds]
	}
	return result
}
