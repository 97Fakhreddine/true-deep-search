package log

import (
	stdlog "log"
	"os"
)

var logger *stdlog.Logger

func init() {
	logger = stdlog.New(os.Stdout, "[hybridsearch] ", stdlog.LstdFlags)
}

func Info(msg string, args ...any) {
	logger.Printf("[INFO] "+msg, args...)
}

func Error(msg string, args ...any) {
	logger.Printf("[ERROR] "+msg, args...)
}

func Debug(msg string, args ...any) {
	logger.Printf("[DEBUG] "+msg, args...)
}
