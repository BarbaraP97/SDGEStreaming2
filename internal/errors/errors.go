// internal/errors/errors.go
package errors

import "fmt"

// AppError represents a standardized application error.
// It is used to provide structured and consistent error handling.
type AppError struct {
	Code    string // Unique code for type of error
	Message string // Human-readable message
	Err     error  // Original error (optional)
}

// Error implements the Go error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Wrap wraps an existing error into an AppError.
func Wrap(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// New creates a new AppError without an underlying error.
func New(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

//
// Common predefined error constructors
//

// ErrNotFound represents a "resource not found" scenario.
func ErrNotFound(resource string) *AppError {
	return New("NOT_FOUND", fmt.Sprintf("%s no encontrado", resource))
}

// ErrInvalidInput represents invalid user input.
func ErrInvalidInput(field string) *AppError {
	return New("INVALID_INPUT", fmt.Sprintf("El campo '%s' no es válido", field))
}

// ErrUnauthorized indicates a missing or invalid auth.
func ErrUnauthorized() *AppError {
	return New("UNAUTHORIZED", "No autorizado")
}

// ErrForbidden indicates lack of permissions.
func ErrForbidden() *AppError {
	return New("FORBIDDEN", "Acceso denegado")
}

// ErrConflict indicates duplicated values (ej: email ya existe).
func ErrConflict(msg string) *AppError {
	return New("CONFLICT", msg)
}

// ErrInternal indicates unexpected server or DB crashes.
func ErrInternal(err error) *AppError {
	return Wrap("INTERNAL_ERROR", "Error interno del servidor", err)
}

// ErrDatabase indicates errors originating from SQL/DB.
func ErrDatabase(err error) *AppError {
	return Wrap("DATABASE_ERROR", "Error de base de datos", err)
}
