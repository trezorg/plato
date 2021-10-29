package requests

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/trezorg/plato/pkg/logger"
)

func (r *Requester) makeRequest(ctx context.Context, method string, url string, client HttpClient, body []byte, query string, debug bool) ([]byte, error) {
	var response []byte
	var err error
	f := func(attempts uint) error {
		var respErr error
		if response, respErr = request(ctx, method, url, client, body, query, debug); respErr != nil {
			if isRetryableError(respErr) {
				// retry, return error
				logger.Infof("Got HTTP error: %v. Retrying. Attempts: %d", respErr, attempts)
				return respErr
			}
			err = respErr
		}
		return nil
	}

	if retryErr := makeRetry(ctx, time.Duration(r.maxTimeout)*time.Second, time.Duration(r.factor)*time.Millisecond, r.jitter, r.attempts, f); retryErr != nil {
		return nil, retryErr
	}

	if err != nil {
		return nil, err
	}
	return response, nil

}

type APIOption func(r *Requester)

func WithRetryConfig(config RetryConfig) APIOption {
	return func(r *Requester) {
		r.RetryConfig = config
	}
}

func WithNewTransport(config TimeoutConfig) APIOption {
	return func(r *Requester) {
		transport, err := newTransport(config)
		if err != nil {
			panic(err)
		}
		client := &http.Client{
			Transport: transport,
			Timeout:   config.requestTimeout,
		}
		r.client = client
	}
}

// Requester structure represents http requests
type Requester struct {
	RetryConfig
	debug  bool
	client HttpClient
}

type Result struct {
	Size  uint
	URL   string
	Error error
}

func (res Result) String() string {
	return fmt.Sprintf("URL: %s. Size: %d, Error: %v", res.URL, res.Size, res.Error)
}

func (r *Requester) Request(ctx context.Context, url string) Result {
	response, err := r.makeRequest(ctx, http.MethodGet, url, r.client, nil, "", r.debug)
	return Result{Size: uint(len(response)), Error: err, URL: url}
}

func (r *Requester) toJob(url string) job {
	return func(ctx context.Context) Result {
		return r.Request(ctx, url)
	}
}

func (r *Requester) toJobs(urls ...string) <-chan job {
	out := make(chan job, len(urls))
	for _, url := range urls {
		out <- r.toJob(url)
	}
	return out
}

func New(timeoutConfig TimeoutConfig, retryConfig RetryConfig, debug bool) (*Requester, error) {
	if err := timeoutConfig.check(); err != nil {
		return nil, err
	}
	if err := retryConfig.check(); err != nil {
		return nil, err
	}
	transport, err := newTransport(timeoutConfig)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeoutConfig.requestTimeout,
	}

	return &Requester{
		RetryConfig: retryConfig,
		client:      client,
		debug:       debug,
	}, nil
}

func Default() (*Requester, error) {
	return New(NewDefaultTimeoutConfig(), NewDefaultRetryConfig(), false)
}

func (r *Requester) Process(ctx context.Context, urls ...string) <-chan Result {
	jobs := r.toJobs(urls...)
	size := uint(min(len(urls), maxPoolSize))
	p := pool{size: size, jobs: jobs}
	return p.start(ctx)
}
