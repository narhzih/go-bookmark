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

/*
TestEmailSignUpFlow tests the email signup flow.
--------------------
# Tested endpoints:
---| /v1/sign-up
---| /v1/verify-account
*/
func TestEmailSignUpFlow(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	t.Run("email sign-up flow", func(t *testing.T) {
		var verificationToken string

		t.Run("/v1/sign-up", func(t *testing.T) {
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

		t.Run("/v1/verify-account", func(t *testing.T) {
			reqUrl := fmt.Sprintf("/v1/verify-account/%v", verificationToken)
			req := httptest.NewRequest(http.MethodPost, reqUrl, nil)
			res := executeRequest(req)
			checkResponseCode(t, http.StatusOK, res.Code)
		})
	})

}

/*
TestUserLoginFlow tests the flow involved in the login process.
--------------------
# Tested endpoints:
---| /v1/login
*/
func TestUserLoginFlow(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	t.Run("login flow", func(t *testing.T) {
		t.Run("/v1/login - success", func(t *testing.T) {
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

/*
TestForgotPassword tests the forgot password flow.
--------------------
# Tested endpoints:
---| /v1/forgot-password
---| /v1/verify-reset-token/:token
---| /v1/reset-password/:token
*/
func TestForgotPassword(t *testing.T) {
	if testing.Short() {
		t.Skip(skipMessage)
	}

	t.Run("forgot password flow", func(t *testing.T) {
		var resetToken string
		t.Run("/v1/forgot-password", func(t *testing.T) {
			resetResData := struct {
				Message string `json:"message"`
				Token   string `json:"token"`
			}{}
			t.Run("/v1/forgot-password | invalid email", func(t *testing.T) {
				resetReqBody := []byte(`{"email": "invaliduser@gmail.com"}`)
				resetReq, err := http.NewRequest(http.MethodPost, "/v1/forgot-password", bytes.NewBuffer(resetReqBody))
				res := executeRequest(resetReq)
				checkResponseCode(t, http.StatusBadRequest, res.Code)

				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("could not read response body %s", err)
				}
				err = json.Unmarshal(resBody, &resetResData)
				if err != nil {
					t.Fatalf("could not unmarshal json body %s", err)
				}
				assert.Equal(t, resetResData.Message, "email does not match any account in our record")
			})

			t.Run("/v1/forgot-password  | valid email", func(t *testing.T) {
				resetReqBody := []byte(`{"email": "user1@gmail.com"}`)
				resetReq, err := http.NewRequest(http.MethodPost, "/v1/forgot-password", bytes.NewBuffer(resetReqBody))
				res := executeRequest(resetReq)
				checkResponseCode(t, http.StatusOK, res.Code)

				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("could not read response body %s", err)
				}
				err = json.Unmarshal(resBody, &resetResData)
				if err != nil {
					t.Fatalf("could not unmarshal json body %s", err)
				}
				resetToken = resetResData.Token
			})
		})

		if resetToken == "" {
			t.Fatal("could not retrieve reset token")
		}

		var verificationSuccessful bool
		t.Run("/v1/verify-reset-token", func(t *testing.T) {
			resetReqData := struct {
				Message string `json:"message"`
				Data    struct {
					User struct {
						ID       int    `json:"id"`
						Username string `json:"username"`
						Email    string `json:"email"`
					} `json:"user"`
					Token string `json:"token"`
				} `json:"data"`
			}{}
			reqUrl := fmt.Sprintf("/v1/verify-reset-token/%s", resetToken)
			req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
			if err != nil {
				t.Fatalf("could not build request %s", err)
			}
			res := executeRequest(req)
			checkResponseCode(t, http.StatusOK, res.Code)
			verificationSuccessful = true
			resBody, err := json.Marshal(res.Body)
			if err != nil {
				t.Fatalf("could not marshal response body %s", err)
			}

			err = json.Unmarshal(resBody, &resetReqData)
			if err != nil {
				t.Fatalf("could not unmarshal response body %s", err)
			}
		})

		if verificationSuccessful != true {
			t.Fatal("could not verify reset token")
		}

		t.Run("/v1/reset-password/", func(t *testing.T) {
			reqUrl := fmt.Sprintf("/v1/reset-password/%s", resetToken)
			reqBody := []byte(`{"password": "Password123"}`)
			req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("could not build request %s", err)
			}
			res := executeRequest(req)
			checkResponseCode(t, http.StatusOK, res.Code)
		})
	})
}

/*
TestTwitterConnectionFlow tests the twitter authentication flow.
--------------------
# Tested endpoints:
---| /v1/twitter/connect-account
---| /v1/twitter/connected-account
---| /v1/twitter/disconnect-account
*/
func TestTwitterConnectionFlow(t *testing.T) {
	t.Run("twitter auth flow", func(t *testing.T) {
		t.Run("/v1/auth/twitter/connect-account", func(t *testing.T) {
			// this will eventually try to connect a user's twitter account
		})
	})
}
