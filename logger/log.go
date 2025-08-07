package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// LogLevel type for custom log levels
type LogLevel int

const (
	FATAL LogLevel = iota
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

var levelLabels = map[LogLevel]string{
	FATAL: colorRed + "FATAL" + colorReset,
	ERROR: colorRed + "ERROR" + colorReset,
	WARN:  colorYellow + "WARN" + colorReset,
	INFO:  colorGreen + "INFO" + colorReset, // Green color for INFO
	DEBUG: colorCyan + "DEBUG" + colorReset,
	TRACE: colorPurple + "TRACE" + colorReset,
}

// Logger struct wraps the standard logger
type Logger struct {
	mu      sync.Mutex
	level   LogLevel
	logger  *log.Logger
	logFile io.Writer
}

// NewLogger creates a new Logger with specified minimum level and output (stdout or file)
func NewLogger(minLevel LogLevel, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}
	return &Logger{
		level:   minLevel,
		logger:  log.New(output, "", 0), // disable default flags, we add custom timestamp
		logFile: output,
	}
}

// SetLevel sets the minimum log level for the logger
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// logMessage formats and logs the message if above set log level
func (l *Logger) logMessage(level LogLevel, format string, v ...interface{}) {
	if level > l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	prefix := fmt.Sprintf("%s [%s] ", timestamp, levelLabels[level])
	msg := fmt.Sprintf(format, v...)
	l.logger.Println(prefix + msg)

	if level == FATAL {
		os.Exit(1)
	}
}

// Public logging functions:

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.logMessage(FATAL, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.logMessage(ERROR, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.logMessage(WARN, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.logMessage(INFO, format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.logMessage(DEBUG, format, v...)
}

func (l *Logger) Trace(format string, v ...interface{}) {
	l.logMessage(TRACE, format, v...)
}
