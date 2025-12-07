package constants

// HTTP Status mapping to standard HTTP codes
var HTTPStatusMap = map[int]int{
	// Success -> 200
	StatusSuccess:   200,
	StatusCreated:   201,
	StatusUpdated:   200,
	StatusDeleted:   200,
	StatusNoContent: 204,

	// Client Errors -> 400
	StatusBadRequest:      400,
	StatusValidationError: 400,
	StatusInvalidJSON:     400,
	StatusMissingField:    400,
	StatusInvalidField:    400,

	// Auth Errors -> 401/403
	StatusUnauthorized:        401,
	StatusInvalidCredentials:  401,
	StatusTokenExpired:        401,
	StatusTokenInvalid:        401,
	StatusRefreshTokenExpired: 401,
	StatusForbidden:           403,
	StatusAccountDisabled:     403,
	StatusEmailNotVerified:    403,
	StatusEmailAlreadyExists:  409,
	StatusOAuthError:          401,

	// Scan Errors -> 400/404
	StatusScanNotFound:      404,
	StatusScanInProgress:    202,
	StatusScanFailed:        400,
	StatusInvalidImage:      400,
	StatusImageTooLarge:     413,
	StatusOCRFailed:         422,
	StatusNutritionNotFound: 422,

	// Product Errors -> 400/404
	StatusProductNotFound:  404,
	StatusBarcodeInvalid:   400,
	StatusBarcodeNotFound:  404,
	StatusProductDataError: 502,

	// User Errors -> 404
	StatusUserNotFound:      404,
	StatusProfileIncomplete: 400,

	// Compare Errors -> 400
	StatusCompareMinTwo:      400,
	StatusCompareMaxExceeded: 400,

	// Rate Limit -> 429
	StatusRateLimitExceeded: 429,
	StatusTooManyRequests:   429,

	// Not Found -> 404
	StatusNotFound:      404,
	StatusRouteNotFound: 404,

	// Server Errors -> 500
	StatusInternalError:    500,
	StatusDatabaseError:    500,
	StatusStorageError:     500,
	StatusExternalAPIError: 502,
	StatusOCRServiceError:  503,
	StatusConfigError:      500,
}

// GetHTTPStatus returns the HTTP status code for an API status code
func GetHTTPStatus(code int) int {
	if status, ok := HTTPStatusMap[code]; ok {
		return status
	}
	return 500
}
