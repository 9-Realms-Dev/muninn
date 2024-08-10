package util

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
)

var (
	Logger *log.Logger
	once   sync.Once
)

func setConfigDirectory() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	munninDir := filepath.Join(configDir, "munnin")
	if err := os.MkdirAll(munninDir, 0755); err != nil {
		return "", err
	}

	munninLogFile := fmt.Sprintf("%s/debug.log", munninDir)

	return munninLogFile, nil
}

func init() {
	once.Do(func() {
		munninLogFile, err := setConfigDirectory()
		if err != nil {
			log.Fatal("Failed to create munnin config directory")
		}

		logFile, err := os.OpenFile(munninLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("Failed to open log file:", err)
		}

		Logger = log.New(logFile)
		Logger.SetReportTimestamp(true)
		Logger.SetReportCaller(true)
		Logger.SetLevel(log.DebugLevel)
	})
}
