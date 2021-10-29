package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	"github.com/trezorg/plato/pkg/logger"
)

const (
	DefaultRetryAttempts  = uint(10)
	DefaultMaxTimeout     = uint(10)
	DefaultRetryJitter    = 0.1
	DefaultRetryFactor    = uint(100)
	DefaultConnectTimeout = uint(5)
	DefaultRequestTimeout = uint(10)
)

type RetryConfig struct {
	attempts uint
	jitter   float64
	// factor in milliseconds
	factor     uint
	maxTimeout uint
}

func (r RetryConfig) check() error {
	if r.attempts == 0 || r.jitter <= 0 || r.factor == 0 {
		return fmt.Errorf("requires all parameters: attempts, jitter, factor")

	}
	return nil
}

func NewRetryConfig(attempts uint, jitter float64, factor uint, maxTimeout uint) RetryConfig {
	return RetryConfig{attempts: attempts, jitter: jitter, factor: factor, maxTimeout: maxTimeout}
}

func NewDefaultRetryConfig() RetryConfig {
	return NewRetryConfig(DefaultRetryAttempts, DefaultRetryJitter, DefaultRetryFactor, DefaultMaxTimeout)
}

type TimeoutConfig struct {
	connectTimeout time.Duration
	requestTimeout time.Duration
}

func NewTimeoutConfig(connectTimeout time.Duration, requestTimeout time.Duration) TimeoutConfig {
	return TimeoutConfig{connectTimeout: connectTimeout, requestTimeout: requestTimeout}
}

func NewDefaultTimeoutConfig() TimeoutConfig {
	return NewTimeoutConfig(time.Duration(DefaultConnectTimeout)*time.Second, time.Duration(DefaultRequestTimeout)*time.Second)
}

func (config TimeoutConfig) check() error {
	if config.connectTimeout == 0 || config.requestTimeout == 0 {
		return fmt.Errorf("requires all parameters: connectTimeout, requestTimeout")
	}
	return nil
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func newTransport(timeoutConfig TimeoutConfig) (*http.Transport, error) {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeoutConfig.connectTimeout,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}, nil
}

func debugRequest(request *http.Request) {
	dump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		logger.Error(err)
	} else {
		logger.Debug(string(dump))
	}
}

func debugResponse(response *http.Response) {
	dump, err := httputil.DumpResponse(response, true)
	if err != nil {
		logger.Error(err)
	} else {
		logger.Debug(string(dump))
	}
}

func request(ctx context.Context, method string, url string, client HttpClient, body []byte, query string, debug bool) ([]byte, error) {
	var requestBody io.Reader = nil
	if len(body) > 0 {
		requestBody = bytes.NewBuffer(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		return nil, RequestError{err: err}
	}
	if len(query) > 0 {
		req.URL.RawQuery = query
	}
	if len(body) > 0 {
		req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}

	for key, values := range agentHeaders {
		for _, header := range values {
			req.Header.Add(key, header)
		}
	}

	if debug {
		debugRequest(req)
	}
	resp, err := client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return nil, NewTimeoutError(err)
		}
		return nil, err
	}
	if debug {
		debugResponse(resp)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}()
	responseBody, err := readBody(resp)
	if err != nil {
		return nil, NewReadBodyError(resp.StatusCode, err)
	}
	if resp.StatusCode >= 400 {
		return nil, NewHTTPError(resp.StatusCode, fmt.Errorf("Http error. Status: %d: %v. Content: %s", resp.StatusCode, resp.Status, responseBody))
	}
	return responseBody, nil
}
