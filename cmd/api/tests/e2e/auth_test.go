package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotest.tools/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmailSignUp(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	t.Run("v1/sign-up", func(t *testing.T) {
		newTestDb(t)
		var verificationToken string
		t.Run("register a new account", func(t *testing.T) {
			signUpRes := struct {
				Message string `json:"message"`
				Data    struct {
					VToken string `json:"v_token"`
				} `json:"data"`
			}{}
			reqBody := []byte(`{"username": "user5", "email": "user5@gmail.com", "password": "Password123", "profile_name": "user5"}`)
			req, err := http.NewRequest(http.MethodPost, "/v1/sign-up", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Errorf("could not create request %s", err)
			}
			res := executeRequest(req)
			checkResponseCode(t, http.StatusCreated, res.Code)
			// further check for the data returned and see if it matches
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf(fmt.Sprintf("could not read response body"))
			}
			err = json.Unmarshal(resBody, &signUpRes)
			if err != nil {
				t.Errorf(fmt.Sprintf(err.Error()))
			}
			assert.Assert(t, signUpRes.Data.VToken != "")
			verificationToken = signUpRes.Data.VToken
		})

		t.Run("account verification after sign-up", func(t *testing.T) {
			reqUrl := fmt.Sprintf("/v1/verify-account/%v", verificationToken)
			req := httptest.NewRequest(http.MethodPost, reqUrl, nil)
			res := executeRequest(req)
			checkResponseCode(t, http.StatusOK, res.Code)
		})
	})

}
