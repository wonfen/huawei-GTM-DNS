package hwdns

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Message(t *testing.T) {
	err := &APIError{StatusCode: 500, Body: "internal oops"}
	assert.Contains(t, err.Error(), "500")
	assert.Contains(t, err.Error(), "internal oops")
}

func TestAPIError_IsRetryable(t *testing.T) {
	cases := []struct {
		code int
		want bool
	}{
		{500, true}, {502, true}, {503, true}, {504, true},
		{429, true},
		{400, false}, {401, false}, {403, false}, {404, false}, {409, false},
		{200, false},
	}
	for _, c := range cases {
		err := &APIError{StatusCode: c.code}
		assert.Equal(t, c.want, err.IsRetryable(), "code %d", c.code)
	}
}

func TestAPIError_UnwrapsViaErrorsAs(t *testing.T) {
	original := &APIError{StatusCode: 503, Body: "temp unavailable"}

	wrappedW := errorsWrap(original)
	var api *APIError
	assert.True(t, errors.As(wrappedW, &api))
	assert.Equal(t, 503, api.StatusCode)
}

func errorsWrap(err error) error {
	return &wrapped{err}
}

type wrapped struct{ inner error }

func (w *wrapped) Error() string { return "wrapped: " + w.inner.Error() }
func (w *wrapped) Unwrap() error { return w.inner }
