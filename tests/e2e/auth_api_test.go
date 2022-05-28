package e2e

import (
	"bytes"
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	t.Run("/signup", func(t *testing.T) {
		t.Run("/successful signup with email", func(t *testing.T) {
			body := []byte(`{"username": "test_userr", "email": "test@exampler.com", "password": "test", "profile_name": "test_profiler"}`)
			req, _ := http.NewRequest(http.MethodPost, "/v1/sign-up", bytes.NewBuffer(body))
			res := executeRequest(req)
			checkResponseCode(t, http.StatusOK, res.Code)
			logger.Info().Msg("All well and Good")
		})
	})
}
