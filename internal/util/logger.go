package util

import (
	"os"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func init() {
	Logger = log.New(os.Stderr)
	Logger.SetReportTimestamp(true)
	Logger.SetReportCaller(true)
	Logger.SetLevel(log.DebugLevel)
}
