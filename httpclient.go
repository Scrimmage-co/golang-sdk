package scrimmage

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

type RetryClient struct {
	client     *http.Client
	maxRetries int
	baseDelay  time.Duration
	maxBackoff time.Duration
}

func NewRetryClient(maxRetries int, baseDelay, maxBackoff time.Duration) *RetryClient {
	return &RetryClient{
		client:     &http.Client{},
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
		maxBackoff: maxBackoff,
	}
}

func (c *RetryClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= c.maxRetries; i++ {
		resp, err = c.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		if resp != nil && resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return resp, nil
		}

		fmt.Printf("Retry %d/%d - Error: %v\n", i+1, c.maxRetries, err)

		delay := c.baseDelay * time.Duration(math.Pow(2, float64(i)))
		if delay > c.maxBackoff {
			delay = c.maxBackoff
		}

		time.Sleep(delay)
	}

	return resp, err
}
