package hwdns

import "fmt"

// APIError represents a non-2xx response from the Huawei DNS API.
// Callers can use errors.As to extract it and read StatusCode for retry decisions.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("huawei api error: status %d, body: %s", e.StatusCode, e.Body)
}

// IsRetryable reports whether the HTTP status is one that callers should
// retry. 5xx and 429 are retryable; other 4xx and non-error codes are not.
func (e *APIError) IsRetryable() bool {
	if e.StatusCode == 429 {
		return true
	}
	return e.StatusCode >= 500 && e.StatusCode < 600
}
