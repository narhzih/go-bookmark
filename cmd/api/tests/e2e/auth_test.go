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
				t.Fatalf("could not create request %s", err)
			}
			res := executeRequest(req)
			checkResponseCode(t, http.StatusCreated, res.Code)
			// further check for the data returned and see if it matches
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal("could not parse response body")
			}
			err = json.Unmarshal(resBody, &signUpRes)
			if err != nil {
				t.Fatal(fmt.Sprintf(err.Error()))
			}
			verificationToken = signUpRes.Data.VToken
		})

		if verificationToken == "" {
			t.Fatal("the verification token is empty, cannot proceed to test other things...")
		}

		t.Run("account verification after sign-up", func(t *testing.T) {
			reqUrl := fmt.Sprintf("/v1/verify-account/%v", verificationToken)
			req := httptest.NewRequest(http.MethodPost, reqUrl, nil)
			res := executeRequest(req)
			checkResponseCode(t, http.StatusOK, res.Code)
		})
	})

}

func TestUserLogin(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	t.Run("/v1/sign-in", func(t *testing.T) {
		newTestDb(t)
		t.Run("successful login", func(t *testing.T) {
			loginResData := struct {
				Message string `json:"message"`
				Data    struct {
					Token        string `json:"token"`
					RefreshToken string `json:"refresh_token"`
					ExpiresAt    string `json:"expires_at"`
					User         struct {
						Id       int    `json:"id"`
						Username string `json:"username"`
						Email    string `json:"email"`
					} `json:"user"`
				} `json:"data"`
			}{}
			loginReqBody := []byte(`{"email": "user1@gmail.com", "password": "password"}`)
			loginReq, err := http.NewRequest(http.MethodPost, "/v1/sign-in", bytes.NewBuffer(loginReqBody))
			if err != nil {
				t.Fatalf("could not build request %s", err)
			}
			loginRes := executeRequest(loginReq)
			loginResBody, err := io.ReadAll(loginRes.Body)
			if err != nil {
				t.Fatalf("could not read login response body %s", err)
			}
			err = json.Unmarshal(loginResBody, &loginResData)
			if err != nil {
				t.Fatalf("could not unmarshal login response body: %s", err)
			}
			t.Log(loginResData.Message)
			checkResponseCode(t, http.StatusOK, loginRes.Code)

			// properly inspect the response we got
			// assert that the returned user credentials matches
			assert.Equal(t, loginResData.Data.User.Username, "user1")
			assert.Equal(t, loginResData.Data.User.Email, "user1@gmail.com")
			assert.Equal(t, loginResData.Data.User.Id, 1)

			// further make sure that we can make a valid authenticated request
			// with the jwt token returned from the request
			authTestReq, err := http.NewRequest(http.MethodGet, "/v1/user/profile", nil)
			if err != nil {
				t.Fatalf("could not build request %s", err)
			}
			authTestReq.Header.Set("Authorization", "Bearer "+loginResData.Data.Token)
			authTestRes := executeRequest(authTestReq)
			checkResponseCode(t, http.StatusOK, authTestRes.Code)
		})
	})
}
