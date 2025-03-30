package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimitMiddleware(t *testing.T) {

	t.Run("limit is not exceeded", func(t *testing.T) {
		application := application{
			config: config{
				limiter: struct {
					rps     float64
					burst   int
					enabled bool
				}{rps: 2, burst: 4, enabled: true},
			},
		}

		responseRecorder := httptest.NewRecorder()

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, []byte("OK"))
		})

		for range 2 {
			request := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
			application.rateLimit(next).ServeHTTP(responseRecorder, request)
		}

		if responseRecorder.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, responseRecorder.Code)
		}
	})

	t.Run("limit is exceeded", func(t *testing.T) {
		application := application{
			config: config{
				limiter: struct {
					rps     float64
					burst   int
					enabled bool
				}{rps: 2, burst: 4, enabled: true},
			},
		}

		next := http.HandlerFunc(application.healthCheckHandler)

		handler := application.rateLimit(next)
		var expected int
		for range 5 {
			// need to put inside loop to create a new recorder for each request. otherwise the recorder will be closed after the first request and return 200 status code for subsequence requests.
			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
			handler.ServeHTTP(responseRecorder, request)
			expected = responseRecorder.Code
		}

		if expected != http.StatusTooManyRequests {
			t.Errorf("expected status code %d but got %d", http.StatusTooManyRequests, expected)
		}

	})

}
