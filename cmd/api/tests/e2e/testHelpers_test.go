package e2e

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/services"
	"github.com/mypipeapp/mypipeapi/cmd/api/services/mailer"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

const (
	skipMessage = "postgres: skipping integration test"
)

// execSqlScript is a helper function to execute SQL commands in the file at the given scriptPath.
func execSqlScript(db *sql.DB, scriptPath string) {
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}
}

func initJWTConfig() (services.JWTConfig, error) {
	var expiresIn int
	var key string
	var err error
	var cfg services.JWTConfig

	expiresIn, err = strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		return cfg, err
	}

	key = os.Getenv("JWT_SECRET")

	// enforce minimum length for JWT secret
	if len(key) < 64 {
		return cfg, fmt.Errorf("JWT_SECRET too short")
	}

	cfg.ExpiresIn = expiresIn
	cfg.Key = key
	cfg.Algo = jwt.SigningMethodHS256

	return cfg, nil
}

func initMailer() *mailer.Mailer {
	var mailerP *mailer.Mailer
	var config mailer.MailConfig

	// I'll probably still have to find a way to use a
	// test mailer server
	config.Password = os.Getenv("MAIL_PASSWORD")
	config.Username = os.Getenv("MAIL_USERNAME")
	config.SmtpHost = os.Getenv("MAIL_HOST")
	config.SmtpPort = os.Getenv("MAIL_PORT")
	config.MailFrom = "My Pipe Desk <desk@mypipe.app>"
	mailerP = mailer.NewMailer(config)
	return mailerP
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.Handler.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func createApplicationInstance() internal.Application {
	jwtConfig, err := initJWTConfig()
	if err != nil {
		logger.Err(err).Msg("jwt config")
	}
	mailerP := initMailer()
	repositories := repository.Repositories{
		User:                postgres.NewUserActions(db, logger),
		Pipe:                postgres.NewPipeActions(db, logger),
		PipeShare:           postgres.NewPipeShareActions(db, logger),
		Bookmark:            postgres.NewBookmarkActions(db, logger),
		Notification:        postgres.NewNotificationActions(db, logger),
		AccountVerification: postgres.NewAccountVerificationActions(db, logger),
		PasswordReset:       postgres.NewPasswordResetActions(db, logger),
		Tag:                 postgres.NewTagActions(db, logger),
		Search:              postgres.NewSearchActions(db, logger),
	}

	appInstance := internal.Application{
		Repositories: repositories,
		Services: services.Services{
			Repositories: repositories,
			Logger:       logger,
			JWTConfig:    jwtConfig,
			Mailer:       mailerP,
		},
		Logger: logger,
	}
	return appInstance
}

func createGlobalUserAndLogin() {
	// Create the account
	signUpRes := struct {
		Message string `json:"message"`
		Data    struct {
			VToken string `json:"v_token"`
		} `json:"data"`
	}{}
	reqBody := []byte(`{"username": "dummy", "email": "dummy@gmail.com", "password": "password", "profile_name": "dummy pn"}`)
	req, _ := http.NewRequest(http.MethodPost, "/v1/sign-up", bytes.NewBuffer(reqBody))

	res := executeRequest(req)
	resBody, _ := io.ReadAll(res.Body)
	err := json.Unmarshal(resBody, &signUpRes)
	if err != nil {
		log.Fatal(fmt.Sprintf(err.Error()))
	}
	verificationToken := signUpRes.Data.VToken

	// verify account
	reqUrl := fmt.Sprintf("/v1/verify-account/%v", verificationToken)
	req, _ = http.NewRequest(http.MethodPost, reqUrl, nil)
	res = executeRequest(req)

	// log user in
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
	loginReqBody := []byte(`{"email": "dummy@gmail.com", "password": "password"}`)
	loginReq, _ := http.NewRequest(http.MethodPost, "/v1/sign-in", bytes.NewBuffer(loginReqBody))

	loginRes := executeRequest(loginReq)
	loginResBody, err := io.ReadAll(loginRes.Body)
	if err != nil {
		log.Fatalf("could not read login response body %s", err)
	}
	err = json.Unmarshal(loginResBody, &loginResData)
	if err != nil {
		log.Fatalf("could not unmarshal login response body: %s", err)
	}

	globalAccessToken = loginResData.Data.Token
	globalUserEmail = loginResData.Data.User.Email
	globalUserID = loginResData.Data.User.Id
}
