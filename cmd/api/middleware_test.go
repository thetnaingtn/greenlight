package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimitMiddleware(t *testing.T) {

	application := application{
		config: config{
			limiter: struct {
				rps     float64
				burst   int
				enabled bool
			}{rps: 2, burst: 4, enabled: true},
		},
	}

	cases := []struct {
		name         string
		expect       int
		noOfRequests int
	}{
		{
			"limit is not exceeded",
			http.StatusOK,
			2,
		},
		{
			"limit is exceeded",
			http.StatusTooManyRequests,
			5,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			next := http.HandlerFunc(application.healthCheckHandler)
			handler := application.rateLimit(next)

			var actual int
			for range c.noOfRequests {
				request := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
				responseRecorder := httptest.NewRecorder()
				handler.ServeHTTP(responseRecorder, request)
				actual = responseRecorder.Code
			}

			if actual != c.expect {
				t.Errorf("expected status code %d but got %d", c.expect, actual)
			}
		})
	}
}
