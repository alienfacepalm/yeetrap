package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeAuth      ErrorType = "authentication"
	ErrorTypeConfig    ErrorType = "configuration"
	ErrorTypeNetwork   ErrorType = "network"
	ErrorTypeFile      ErrorType = "filesystem"
	ErrorTypeAPI       ErrorType = "api"
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeExternal  ErrorType = "external"
)

// YeeTrapError represents a custom error with additional context
type YeeTrapError struct {
	Type    ErrorType
	Message string
	Details string
	Cause   error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *YeeTrapError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *YeeTrapError) Unwrap() error {
	return e.Cause
}

// WithContext adds context to the error
func (e *YeeTrapError) WithContext(key string, value interface{}) *YeeTrapError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithDetails adds details to the error
func (e *YeeTrapError) WithDetails(details string) *YeeTrapError {
	e.Details = details
	return e
}

// New creates a new YeeTrapError
func New(errorType ErrorType, message string) *YeeTrapError {
	return &YeeTrapError{
		Type:    errorType,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errorType ErrorType, message string) *YeeTrapError {
	return &YeeTrapError{
		Type:    errorType,
		Message: message,
		Cause:   err,
		Context: make(map[string]interface{}),
	}
}

// Predefined error constructors
func NewAuthError(message string) *YeeTrapError {
	return New(ErrorTypeAuth, message)
}

func NewConfigError(message string) *YeeTrapError {
	return New(ErrorTypeConfig, message)
}

func NewNetworkError(message string) *YeeTrapError {
	return New(ErrorTypeNetwork, message)
}

func NewFileError(message string) *YeeTrapError {
	return New(ErrorTypeFile, message)
}

func NewAPIError(message string) *YeeTrapError {
	return New(ErrorTypeAPI, message)
}

func NewValidationError(message string) *YeeTrapError {
	return New(ErrorTypeValidation, message)
}

func NewExternalError(message string) *YeeTrapError {
	return New(ErrorTypeExternal, message)
}

// WrapAuth wraps an error as an authentication error
func WrapAuth(err error, message string) *YeeTrapError {
	return Wrap(err, ErrorTypeAuth, message)
}

// WrapConfig wraps an error as a configuration error
func WrapConfig(err error, message string) *YeeTrapError {
	return Wrap(err, ErrorTypeConfig, message)
}

// WrapNetwork wraps an error as a network error
func WrapNetwork(err error, message string) *YeeTrapError {
	return Wrap(err, ErrorTypeNetwork, message)
}

// WrapFile wraps an error as a filesystem error
func WrapFile(err error, message string) *YeeTrapError {
	return Wrap(err, ErrorTypeFile, message)
}

// WrapAPI wraps an error as an API error
func WrapAPI(err error, message string) *YeeTrapError {
	return Wrap(err, ErrorTypeAPI, message)
}

// WrapValidation wraps an error as a validation error
func WrapValidation(err error, message string) *YeeTrapError {
	return Wrap(err, ErrorTypeValidation, message)
}

// WrapExternal wraps an error as an external tool error
func WrapExternal(err error, message string) *YeeTrapError {
	return Wrap(err, ErrorTypeExternal, message)
}

// IsYeeTrapError checks if an error is a YeeTrapError
func IsYeeTrapError(err error) bool {
	_, ok := err.(*YeeTrapError)
	return ok
}

// GetErrorType returns the error type if it's a YeeTrapError
func GetErrorType(err error) ErrorType {
	if ytErr, ok := err.(*YeeTrapError); ok {
		return ytErr.Type
	}
	return ""
}

// FormatError formats an error for user display
func FormatError(err error) string {
	if ytErr, ok := err.(*YeeTrapError); ok {
		var parts []string
		
		// Add main message
		parts = append(parts, ytErr.Message)
		
		// Add details if available
		if ytErr.Details != "" {
			parts = append(parts, ytErr.Details)
		}
		
		// Add context if available
		if len(ytErr.Context) > 0 {
			for key, value := range ytErr.Context {
				parts = append(parts, fmt.Sprintf("%s: %v", key, value))
			}
		}
		
		return strings.Join(parts, " - ")
	}
	
	return err.Error()
}

// GetUserFriendlyMessage returns a user-friendly error message
func GetUserFriendlyMessage(err error) string {
	if ytErr, ok := err.(*YeeTrapError); ok {
		switch ytErr.Type {
		case ErrorTypeAuth:
			return "Authentication failed. Please run 'yeetrap auth' to authenticate."
		case ErrorTypeConfig:
			return "Configuration error. Please run 'yeetrap setup' to configure the application."
		case ErrorTypeNetwork:
			return "Network error. Please check your internet connection and try again."
		case ErrorTypeFile:
			return "File system error. Please check file permissions and disk space."
		case ErrorTypeAPI:
			return "YouTube API error. Please try again later or check your API quota."
		case ErrorTypeValidation:
			return "Invalid input. Please check your command parameters."
		case ErrorTypeExternal:
			return "External tool error. Please ensure yt-dlp is installed and accessible."
		default:
			return ytErr.Message
		}
	}
	
	return err.Error()
}
