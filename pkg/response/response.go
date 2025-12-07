package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Metadata contains additional information about the response
// @Description Response metadata
type Metadata struct {
	Timestamp string `json:"timestamp" example:"2024-12-07T10:00:00Z"`
	RequestID string `json:"request_id,omitempty" example:"abc123"`
	Version   string `json:"version,omitempty" example:"1.0.0"`
}

// Pagination contains pagination information
// @Description Pagination metadata
type Pagination struct {
	Page       int   `json:"page" example:"1"`
	PerPage    int   `json:"per_page" example:"10"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
}

// ErrorDetail contains error information
// @Description Error details
type ErrorDetail struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid request"`
	Field   string `json:"field,omitempty" example:"email"`
}

// SuccessEnvelope is the standard success response wrapper
// @Description Standard success response envelope
type SuccessEnvelope struct {
	Success    bool        `json:"success" example:"true"`
	Data       interface{} `json:"data"`
	Metadata   *Metadata   `json:"metadata,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// ErrorEnvelope is the standard error response wrapper
// @Description Standard error response envelope
type ErrorEnvelope struct {
	Success  bool          `json:"success" example:"false"`
	Error    ErrorDetail   `json:"error"`
	Errors   []ErrorDetail `json:"errors,omitempty"`
	Metadata *Metadata     `json:"metadata,omitempty"`
}

// MessageEnvelope is a simple message response
// @Description Simple message response
type MessageEnvelope struct {
	Success  bool      `json:"success" example:"true"`
	Message  string    `json:"message" example:"Operation completed successfully"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

func newMetadata() *Metadata {
	return &Metadata{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
	}
}

// Success sends a success response with data
func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(SuccessEnvelope{
		Success:  true,
		Data:     data,
		Metadata: newMetadata(),
	})
}

// Created sends a 201 success response with data
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(SuccessEnvelope{
		Success:  true,
		Data:     data,
		Metadata: newMetadata(),
	})
}

// SuccessWithPagination sends a success response with pagination
func SuccessWithPagination(c *fiber.Ctx, data interface{}, pagination Pagination) error {
	return c.Status(fiber.StatusOK).JSON(SuccessEnvelope{
		Success:    true,
		Data:       data,
		Metadata:   newMetadata(),
		Pagination: &pagination,
	})
}

// Message sends a simple message response
func Message(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusOK).JSON(MessageEnvelope{
		Success:  true,
		Message:  message,
		Metadata: newMetadata(),
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(ErrorEnvelope{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
		Metadata: newMetadata(),
	})
}

// ErrorWithField sends an error response with field information
func ErrorWithField(c *fiber.Ctx, code int, message string, field string) error {
	return c.Status(code).JSON(ErrorEnvelope{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Field:   field,
		},
		Metadata: newMetadata(),
	})
}

// ValidationErrors sends multiple validation errors
func ValidationErrors(c *fiber.Ctx, errors []ErrorDetail) error {
	return c.Status(fiber.StatusBadRequest).JSON(ErrorEnvelope{
		Success: false,
		Error: ErrorDetail{
			Code:    fiber.StatusBadRequest,
			Message: "Validation failed",
		},
		Errors:   errors,
		Metadata: newMetadata(),
	})
}

// BadRequest sends a 400 error response
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

// Unauthorized sends a 401 error response
func Unauthorized(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return Error(c, fiber.StatusUnauthorized, message)
}

// Forbidden sends a 403 error response
func Forbidden(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return Error(c, fiber.StatusForbidden, message)
}

// NotFound sends a 404 error response
func NotFound(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return Error(c, fiber.StatusNotFound, message)
}

// InternalError sends a 500 error response
func InternalError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return Error(c, fiber.StatusInternalServerError, message)
}
