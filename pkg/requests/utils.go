package requests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-multierror"
	"github.com/trezorg/plato/pkg/logger"
)

func readBody(resp *http.Response) ([]byte, error) {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Error(err)
		}
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := data
	return res, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ParseURLs(urls ...string) ([]string, error) {
	var result *multierror.Error
	var parsed []string
	for _, u := range urls {
		uri, err := url.Parse(u)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("Malformed url: %s. %w", u, err))
		}
		if uri.Scheme == "" {
			uri.Scheme = "https"
		}
		parsed = append(parsed, uri.String())
	}
	return parsed, result.ErrorOrNil()
}
