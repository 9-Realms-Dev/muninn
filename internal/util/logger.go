package util

import (
	"os"
	"sync"

	"github.com/charmbracelet/log"
)

var (
	Logger *log.Logger
	once   sync.Once
)

func init() {
	once.Do(func() {
		logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("Failed to open log file:", err)
		}

		Logger = log.New(logFile)
		Logger.SetReportTimestamp(true)
		Logger.SetReportCaller(true)
		Logger.SetLevel(log.DebugLevel)
	})
}
