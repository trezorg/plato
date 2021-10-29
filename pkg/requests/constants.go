package requests

import "net/http"

var (
	agentHeaders = http.Header{
		"User-Agent": []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36"},
		"Accept":     []string{"application/json"},
	}
)
