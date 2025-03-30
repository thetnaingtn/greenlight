package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {

	application := application{}

	request := httptest.NewRequest("GET", "/v1/healthcheck", nil)
	responseRecorder := httptest.NewRecorder()

	application.healthCheckHandler(responseRecorder, request)

	result := responseRecorder.Result()

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)

	if err != nil {
		t.Fatal("failed to read response body")
	}

	if result.StatusCode != 200 {
		t.Errorf("expected status code 200 but got %d", result.StatusCode)
	}

	var data envelope
	if err = json.Unmarshal(body, &data); err != nil {
		log.Println(err)
		t.Fatal("failed to parse response body")
	}

	if data["status"] != "available" {
		t.Errorf("expected status to be available but got %s", data["status"])
	}

	systemInfo, ok := data["system_info"].(map[string]any)
	if !ok {
		t.Fatal("failed to parse system_info")
	}

	if systemInfo["environment"] != "development" {
		t.Errorf("expected environment to be development but got %s", systemInfo["environment"])
	}

	if systemInfo["version"] != version {
		t.Errorf("expected version to be %s but got %s", version, systemInfo["version"])
	}

}
