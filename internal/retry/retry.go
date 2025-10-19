package retry

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/AlienFacepalm/YeeTrap/internal/errors"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Multiplier  float64
	Jitter      bool
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   time.Second,
		MaxDelay:    30 * time.Second,
		Multiplier:  2.0,
		Jitter:      true,
	}
}

// NetworkRetryConfig returns a retry configuration for network operations
func NetworkRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   2 * time.Second,
		MaxDelay:    60 * time.Second,
		Multiplier:  2.0,
		Jitter:      true,
	}
}

// APIRetryConfig returns a retry configuration for API operations
func APIRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Second,
		MaxDelay:    10 * time.Second,
		Multiplier:  1.5,
		Jitter:      false,
	}
}

// DownloadRetryConfig returns a retry configuration for download operations
func DownloadRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts: 2,
		BaseDelay:   5 * time.Second,
		MaxDelay:    30 * time.Second,
		Multiplier:  2.0,
		Jitter:      true,
	}
}

// RetryableFunc is a function that can be retried
type RetryableFunc func() error

// RetryableFuncWithContext is a function that can be retried with context
type RetryableFuncWithContext func(ctx context.Context) error

// Retry executes a function with retry logic
func Retry(fn RetryableFunc, config *RetryConfig) error {
	return RetryWithContext(context.Background(), func(ctx context.Context) error {
		return fn()
	}, config)
}

// RetryWithContext executes a function with retry logic and context
func RetryWithContext(ctx context.Context, fn RetryableFuncWithContext, config *RetryConfig) error {
	if config == nil {
		config = DefaultRetryConfig()
	}
	
	var lastErr error
	
	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		// Execute the function
		err := fn(ctx)
		if err == nil {
			return nil // Success
		}
		
		lastErr = err
		
		// Don't retry on the last attempt
		if attempt == config.MaxAttempts {
			break
		}
		
		// Check if error is retryable
		if !IsRetryableError(err) {
			return err
		}
		
		// Calculate delay
		delay := calculateDelay(attempt, config)
		
		// Log retry attempt
		fmt.Printf("Attempt %d failed: %v. Retrying in %v...\n", attempt, err, delay)
		
		// Wait before retry
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	
	// All attempts failed
	return errors.WrapAPI(lastErr, fmt.Sprintf("operation failed after %d attempts", config.MaxAttempts))
}

// IsRetryableError checks if an error should be retried
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check if it's a YeeTrapError
	if ytErr, ok := err.(*errors.YeeTrapError); ok {
		switch ytErr.Type {
		case errors.ErrorTypeNetwork:
			return true
		case errors.ErrorTypeAPI:
			// Check for specific API errors that should be retried
			msg := ytErr.Error()
			if containsAny(msg, []string{"quota", "rate limit", "temporary", "timeout", "connection"}) {
				return true
			}
			return false
		case errors.ErrorTypeExternal:
			// Check for specific external tool errors
			msg := ytErr.Error()
			if containsAny(msg, []string{"timeout", "connection", "network", "temporary"}) {
				return true
			}
			return false
		default:
			return false
		}
	}
	
	// Check for common retryable error patterns
	msg := err.Error()
	retryablePatterns := []string{
		"timeout",
		"connection",
		"network",
		"temporary",
		"rate limit",
		"quota",
		"server error",
		"service unavailable",
	}
	
	return containsAny(msg, retryablePatterns)
}

// calculateDelay calculates the delay for the next retry attempt
func calculateDelay(attempt int, config *RetryConfig) time.Duration {
	// Exponential backoff
	delay := float64(config.BaseDelay) * math.Pow(config.Multiplier, float64(attempt-1))
	
	// Cap at max delay
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}
	
	// Add jitter if enabled
	if config.Jitter {
		// Add up to 25% jitter
		jitter := delay * 0.25 * (0.5 - math.Mod(float64(time.Now().UnixNano()), 1.0))
		delay += jitter
	}
	
	return time.Duration(delay)
}

// containsAny checks if a string contains any of the given substrings
func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

// RetryWithBackoff executes a function with exponential backoff
func RetryWithBackoff(fn RetryableFunc, maxAttempts int) error {
	config := &RetryConfig{
		MaxAttempts: maxAttempts,
		BaseDelay:   time.Second,
		MaxDelay:    30 * time.Second,
		Multiplier:  2.0,
		Jitter:      true,
	}
	return Retry(fn, config)
}

// RetryNetworkOperation retries a network operation with appropriate settings
func RetryNetworkOperation(fn RetryableFunc) error {
	return Retry(fn, NetworkRetryConfig())
}

// RetryAPIOperation retries an API operation with appropriate settings
func RetryAPIOperation(fn RetryableFunc) error {
	return Retry(fn, APIRetryConfig())
}

// RetryDownloadOperation retries a download operation with appropriate settings
func RetryDownloadOperation(fn RetryableFunc) error {
	return Retry(fn, DownloadRetryConfig())
}
