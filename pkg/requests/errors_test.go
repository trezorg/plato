package requests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetryableError(t *testing.T) {
	err := fmt.Errorf("test")

	api400Error := NewHTTPError(400, err)
	apiRedirectError1 := NewHTTPError(301, err)
	apiRedirectError2 := NewHTTPError(302, err)
	timeoutError := NewTimeoutError(err)
	api500Error := NewHTTPError(500, err)
	api502Error := NewHTTPError(502, err)

	assert.False(t, isRetryableError(err))
	assert.False(t, isRetryableError(api400Error))
	assert.False(t, isRetryableError(apiRedirectError1))
	assert.False(t, isRetryableError(apiRedirectError2))
	assert.True(t, isRetryableError(timeoutError))
	assert.True(t, isRetryableError(api500Error))
	assert.True(t, isRetryableError(api502Error))

}
