package system_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"swiftly/backend/internal/pkg/response"
)

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, "Welcome to Swiftly API", nil)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp response.APIResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if !resp.Success {
		t.Errorf("handler returned success false, want true")
	}

	if resp.Message != "Welcome to Swiftly API" {
		t.Errorf("handler returned unexpected message: got %v", resp.Message)
	}
}
