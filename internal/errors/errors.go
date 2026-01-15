package errors

import "fmt"

// AppError interface defines the contract for application errors
type AppError interface {
	Error() string
	IsValidationError() bool
	IsBusinessRuleError() bool
	IsSystemError() bool
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}

func (e ValidationError) IsValidationError() bool {
	return true
}

func (e ValidationError) IsBusinessRuleError() bool {
	return false
}

func (e ValidationError) IsSystemError() bool {
	return false
}

// BusinessRuleError represents business rule violations
type BusinessRuleError struct {
	Rule    string
	Message string
	Context map[string]interface{}
}

func (e BusinessRuleError) Error() string {
	return fmt.Sprintf("business rule violation '%s': %s", e.Rule, e.Message)
}

func (e BusinessRuleError) IsValidationError() bool {
	return false
}

func (e BusinessRuleError) IsBusinessRuleError() bool {
	return true
}

func (e BusinessRuleError) IsSystemError() bool {
	return false
}

// SystemError represents system-level errors
type SystemError struct {
	Operation string
	Err       error
	Message   string
}

func (e SystemError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("system error in '%s': %s (%v)", e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("system error in '%s': %v", e.Operation, e.Err)
}

func (e SystemError) Unwrap() error {
	return e.Err
}

func (e SystemError) IsValidationError() bool {
	return false
}

func (e SystemError) IsBusinessRuleError() bool {
	return false
}

func (e SystemError) IsSystemError() bool {
	return true
}

// Helper functions to create errors
func NewValidationError(field, message string, value interface{}) ValidationError {
	return ValidationError{Field: field, Message: message, Value: value}
}

func NewBusinessRuleError(rule, message string) BusinessRuleError {
	return BusinessRuleError{Rule: rule, Message: message, Context: make(map[string]interface{})}
}

func NewBusinessRuleErrorWithContext(rule, message string, context map[string]interface{}) BusinessRuleError {
	return BusinessRuleError{Rule: rule, Message: message, Context: context}
}

func NewSystemError(operation string, err error) SystemError {
	return SystemError{Operation: operation, Err: err}
}

func NewSystemErrorWithMessage(operation, message string, err error) SystemError {
	return SystemError{Operation: operation, Err: err, Message: message}
}
