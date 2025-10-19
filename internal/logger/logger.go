package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a logger instance
type Logger struct {
	level    LogLevel
	logger   *log.Logger
	file     *os.File
	debug    bool
	verbose  bool
}

var (
	defaultLogger *Logger
	initialized   bool
)

// Init initializes the default logger
func Init(level LogLevel, logFile string, debug, verbose bool) error {
	if initialized {
		return nil
	}
	
	var writer io.Writer = os.Stdout
	var file *os.File
	
	// Set up file logging if specified
	if logFile != "" {
		// Ensure log directory exists
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
		
		// Open log file
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		file = f
		
		// Use both stdout and file if verbose
		if verbose {
			writer = io.MultiWriter(os.Stdout, f)
		} else {
			writer = f
		}
	}
	
	defaultLogger = &Logger{
		level:   level,
		logger:  log.New(writer, "", 0),
		file:    file,
		debug:   debug,
		verbose: verbose,
	}
	
	initialized = true
	return nil
}

// GetLogger returns the default logger
func GetLogger() *Logger {
	if !initialized {
		// Initialize with default settings
		Init(LogLevelInfo, "", false, false)
	}
	return defaultLogger
}

// Close closes the logger and any open files
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetDebug sets debug mode
func (l *Logger) SetDebug(debug bool) {
	l.debug = debug
}

// SetVerbose sets verbose mode
func (l *Logger) SetVerbose(verbose bool) {
	l.verbose = verbose
}

// log formats and logs a message
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	
	// Get caller information for debug logs
	var caller string
	if l.debug && level == LogLevelDebug {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			caller = fmt.Sprintf("[%s:%d] ", filepath.Base(file), line)
		}
	}
	
	// Format timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	// Format message
	message := fmt.Sprintf(format, args...)
	
	// Log the message
	l.logger.Printf("[%s] %s%s %s", timestamp, level.String(), caller, message)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LogLevelWarn, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(LogLevelFatal, format, args...)
	os.Exit(1)
}

// Package-level convenience functions
func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Fatal(format, args...)
}

// LogError logs an error with additional context
func LogError(err error, context string) {
	if err != nil {
		Error("%s: %v", context, err)
	}
}

// LogOperation logs the start and end of an operation
func LogOperation(operation string, fn func() error) error {
	Info("Starting %s", operation)
	start := time.Now()
	
	err := fn()
	
	duration := time.Since(start)
	if err != nil {
		Error("Failed %s after %v: %v", operation, duration, err)
	} else {
		Info("Completed %s in %v", operation, duration)
	}
	
	return err
}

// LogOperationWithResult logs an operation with a result
func LogOperationWithResult(operation string, fn func() (interface{}, error)) (interface{}, error) {
	Info("Starting %s", operation)
	start := time.Now()
	
	result, err := fn()
	
	duration := time.Since(start)
	if err != nil {
		Error("Failed %s after %v: %v", operation, duration, err)
	} else {
		Info("Completed %s in %v", operation, duration)
	}
	
	return result, err
}

// LogAPIRequest logs an API request
func LogAPIRequest(method, url string, statusCode int, duration time.Duration) {
	Info("API %s %s - %d - %v", method, url, statusCode, duration)
}

// LogDownload logs download progress
func LogDownload(videoID, title string, progress float64, speed string) {
	Info("Download %s (%s) - %.1f%% - %s", videoID, title, progress, speed)
}

// LogAuth logs authentication events
func LogAuth(event string, details ...interface{}) {
	Info("Auth %s: %v", event, details)
}

// LogConfig logs configuration events
func LogConfig(event string, details ...interface{}) {
	Info("Config %s: %v", event, details)
}

// LogValidation logs validation events
func LogValidation(event string, details ...interface{}) {
	Debug("Validation %s: %v", event, details)
}

// LogRetry logs retry attempts
func LogRetry(operation string, attempt int, maxAttempts int, err error) {
	Warn("Retry %s (attempt %d/%d): %v", operation, attempt, maxAttempts, err)
}

// LogProgress logs progress updates
func LogProgress(operation string, current, total int, item string) {
	Info("Progress %s: %d/%d - %s", operation, current, total, item)
}

// CleanupLogs removes old log files
func CleanupLogs(logDir string, maxAge time.Duration) error {
	if logDir == "" {
		return nil
	}
	
	files, err := filepath.Glob(filepath.Join(logDir, "*.log"))
	if err != nil {
		return err
	}
	
	cutoff := time.Now().Add(-maxAge)
	
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(file); err != nil {
				Warn("Failed to remove old log file %s: %v", file, err)
			} else {
				Info("Removed old log file: %s", file)
			}
		}
	}
	
	return nil
}
