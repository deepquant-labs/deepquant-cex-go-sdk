package errors

import (
	"fmt"
)

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	// General errors
	ErrUnknown         ErrorCode = "UNKNOWN_ERROR"
	ErrInvalidInput    ErrorCode = "INVALID_INPUT"
	ErrTimeout         ErrorCode = "TIMEOUT"
	ErrRateLimit       ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrNetworkError    ErrorCode = "NETWORK_ERROR"
	ErrInvalidResponse ErrorCode = "INVALID_RESPONSE"

	// Authentication errors
	ErrInvalidAPIKey    ErrorCode = "INVALID_API_KEY" // #nosec G101 -- This is an error code, not a credential
	ErrInvalidSignature ErrorCode = "INVALID_SIGNATURE"
	ErrPermissionDenied ErrorCode = "PERMISSION_DENIED"
	ErrAPIKeyExpired    ErrorCode = "API_KEY_EXPIRED" // #nosec G101 -- This is an error code, not a credential

	// Exchange specific errors
	ErrExchangeNotSupported ErrorCode = "EXCHANGE_NOT_SUPPORTED"
	ErrExchangeUnavailable  ErrorCode = "EXCHANGE_UNAVAILABLE"
	ErrInvalidSymbol        ErrorCode = "INVALID_SYMBOL"
	ErrInsufficientBalance  ErrorCode = "INSUFFICIENT_BALANCE"
	ErrOrderNotFound        ErrorCode = "ORDER_NOT_FOUND"
	ErrInvalidOrderType     ErrorCode = "INVALID_ORDER_TYPE"
	ErrAPIError             ErrorCode = "API_ERROR"

	// Data parsing errors
	ErrJSONParsing      ErrorCode = "JSON_PARSING_ERROR"
	ErrDataParsingError ErrorCode = "DATA_PARSING_ERROR"
	ErrDataFormat       ErrorCode = "INVALID_DATA_FORMAT"
	ErrMissingField     ErrorCode = "MISSING_FIELD"
	ErrInvalidDataType  ErrorCode = "INVALID_DATA_TYPE"
)

// SDKError represents a standardized error in the SDK
type SDKError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Cause   error     `json:"-"`
}

// Error implements the error interface
func (e *SDKError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause error
func (e *SDKError) Unwrap() error {
	return e.Cause
}

// New creates a new SDKError
func New(code ErrorCode, message string) *SDKError {
	return &SDKError{
		Code:    code,
		Message: message,
	}
}

// Newf creates a new SDKError with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *SDKError {
	return &SDKError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap wraps an existing error with SDK error information
func Wrap(code ErrorCode, message string, cause error) *SDKError {
	return &SDKError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Wrapf wraps an existing error with formatted message
func Wrapf(code ErrorCode, cause error, format string, args ...interface{}) *SDKError {
	return &SDKError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
	}
}

// WithDetails adds details to an existing SDKError
func (e *SDKError) WithDetails(details string) *SDKError {
	e.Details = details
	return e
}

// WithDetailsf adds formatted details to an existing SDKError
func (e *SDKError) WithDetailsf(format string, args ...interface{}) *SDKError {
	e.Details = fmt.Sprintf(format, args...)
	return e
}

// IsSDKError checks if an error is an SDKError
func IsSDKError(err error) bool {
	_, ok := err.(*SDKError)
	return ok
}

// GetCode extracts the error code from an error
func GetCode(err error) ErrorCode {
	if sdkErr, ok := err.(*SDKError); ok {
		return sdkErr.Code
	}
	return ErrUnknown
}

// Common error constructors for convenience
func ErrInvalidInputf(format string, args ...interface{}) *SDKError {
	return Newf(ErrInvalidInput, format, args...)
}

func ErrNetworkf(format string, args ...interface{}) *SDKError {
	return Newf(ErrNetworkError, format, args...)
}

func ErrRateLimitf(format string, args ...interface{}) *SDKError {
	return Newf(ErrRateLimit, format, args...)
}

func ErrExchangeNotSupportedf(format string, args ...interface{}) *SDKError {
	return Newf(ErrExchangeNotSupported, format, args...)
}

func ErrJSONParsingf(format string, args ...interface{}) *SDKError {
	return Newf(ErrJSONParsing, format, args...)
}
