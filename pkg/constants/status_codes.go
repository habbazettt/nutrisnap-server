package constants

// API Status Codes
// Format: XXYYZZ
// XX = Category (10=Success, 40=Client Error, 50=Server Error)
// YY = Module (00=General, 01=Auth, 02=Scan, 03=Product, 04=User, 05=Compare)
// ZZ = Specific Error

const (
	// ========== SUCCESS (10XXXX) ==========
	StatusSuccess   = 100000 // General success
	StatusCreated   = 100001 // Resource created
	StatusUpdated   = 100002 // Resource updated
	StatusDeleted   = 100003 // Resource deleted
	StatusNoContent = 100004 // No content to return

	// ========== CLIENT ERRORS - General (40XXXX) ==========
	StatusBadRequest      = 400000 // Bad request
	StatusValidationError = 400001 // Validation failed
	StatusInvalidJSON     = 400002 // Invalid JSON format
	StatusMissingField    = 400003 // Required field missing
	StatusInvalidField    = 400004 // Invalid field value

	// ========== CLIENT ERRORS - Auth (4001XX) ==========
	StatusUnauthorized        = 400100 // Not authenticated
	StatusInvalidCredentials  = 400101 // Invalid email/password
	StatusTokenExpired        = 400102 // JWT token expired
	StatusTokenInvalid        = 400103 // JWT token invalid
	StatusRefreshTokenExpired = 400104 // Refresh token expired
	StatusForbidden           = 400105 // No permission
	StatusAccountDisabled     = 400106 // Account is disabled
	StatusEmailNotVerified    = 400107 // Email not verified
	StatusEmailAlreadyExists  = 400108 // Email already registered
	StatusOAuthError          = 400109 // OAuth authentication error

	// ========== CLIENT ERRORS - Scan (4002XX) ==========
	StatusScanNotFound      = 400200 // Scan not found
	StatusScanInProgress    = 400201 // Scan still processing
	StatusScanFailed        = 400202 // Scan processing failed
	StatusInvalidImage      = 400203 // Invalid image format
	StatusImageTooLarge     = 400204 // Image size too large
	StatusOCRFailed         = 400205 // OCR processing failed
	StatusNutritionNotFound = 400206 // Nutrition data not found

	// ========== CLIENT ERRORS - Product (4003XX) ==========
	StatusProductNotFound  = 400300 // Product not found
	StatusBarcodeInvalid   = 400301 // Invalid barcode format
	StatusBarcodeNotFound  = 400302 // Barcode not in database
	StatusProductDataError = 400303 // Error fetching product data

	// ========== CLIENT ERRORS - User (4004XX) ==========
	StatusUserNotFound      = 400400 // User not found
	StatusProfileIncomplete = 400401 // Profile incomplete

	// ========== CLIENT ERRORS - Compare (4005XX) ==========
	StatusCompareMinTwo      = 400500 // Minimum 2 products required
	StatusCompareMaxExceeded = 400501 // Maximum products exceeded

	// ========== CLIENT ERRORS - Rate Limit (4006XX) ==========
	StatusRateLimitExceeded = 400600 // Rate limit exceeded
	StatusTooManyRequests   = 400601 // Too many requests

	// ========== NOT FOUND (404XXX) ==========
	StatusNotFound      = 404000 // Resource not found
	StatusRouteNotFound = 404001 // Route not found

	// ========== SERVER ERRORS (50XXXX) ==========
	StatusInternalError    = 500000 // Internal server error
	StatusDatabaseError    = 500001 // Database error
	StatusStorageError     = 500002 // Storage (MinIO) error
	StatusExternalAPIError = 500003 // External API error
	StatusOCRServiceError  = 500004 // OCR service error
	StatusConfigError      = 500005 // Configuration error
)

// StatusMessages maps status codes to human-readable messages
var StatusMessages = map[int]string{
	// Success
	StatusSuccess:   "Success",
	StatusCreated:   "Resource created successfully",
	StatusUpdated:   "Resource updated successfully",
	StatusDeleted:   "Resource deleted successfully",
	StatusNoContent: "No content",

	// Client Errors - General
	StatusBadRequest:      "Bad request",
	StatusValidationError: "Validation failed",
	StatusInvalidJSON:     "Invalid JSON format",
	StatusMissingField:    "Required field is missing",
	StatusInvalidField:    "Invalid field value",

	// Client Errors - Auth
	StatusUnauthorized:        "Authentication required",
	StatusInvalidCredentials:  "Invalid email or password",
	StatusTokenExpired:        "Token has expired",
	StatusTokenInvalid:        "Invalid token",
	StatusRefreshTokenExpired: "Refresh token has expired",
	StatusForbidden:           "Access forbidden",
	StatusAccountDisabled:     "Account is disabled",
	StatusEmailNotVerified:    "Email not verified",
	StatusEmailAlreadyExists:  "Email already registered",
	StatusOAuthError:          "OAuth authentication failed",

	// Client Errors - Scan
	StatusScanNotFound:      "Scan not found",
	StatusScanInProgress:    "Scan is still processing",
	StatusScanFailed:        "Scan processing failed",
	StatusInvalidImage:      "Invalid image format",
	StatusImageTooLarge:     "Image size too large",
	StatusOCRFailed:         "OCR processing failed",
	StatusNutritionNotFound: "Nutrition data not found in image",

	// Client Errors - Product
	StatusProductNotFound:  "Product not found",
	StatusBarcodeInvalid:   "Invalid barcode format",
	StatusBarcodeNotFound:  "Barcode not found in database",
	StatusProductDataError: "Error fetching product data",

	// Client Errors - User
	StatusUserNotFound:      "User not found",
	StatusProfileIncomplete: "Profile is incomplete",

	// Client Errors - Compare
	StatusCompareMinTwo:      "Minimum 2 products required for comparison",
	StatusCompareMaxExceeded: "Maximum products for comparison exceeded",

	// Client Errors - Rate Limit
	StatusRateLimitExceeded: "Rate limit exceeded",
	StatusTooManyRequests:   "Too many requests, please try again later",

	// Not Found
	StatusNotFound:      "Resource not found",
	StatusRouteNotFound: "Route not found",

	// Server Errors
	StatusInternalError:    "Internal server error",
	StatusDatabaseError:    "Database error",
	StatusStorageError:     "Storage service error",
	StatusExternalAPIError: "External API error",
	StatusOCRServiceError:  "OCR service error",
	StatusConfigError:      "Configuration error",
}

// GetStatusMessage returns the message for a status code
func GetStatusMessage(code int) string {
	if msg, ok := StatusMessages[code]; ok {
		return msg
	}
	return "Unknown error"
}
