package e2e

import (
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/services"
	"github.com/mypipeapp/mypipeapi/cmd/api/services/mailer"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"github.com/mypipeapp/mypipeapi/db/repository"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

const (
	skipMessage = "postgres: skipping integration test"
)

func newTestDb(t *testing.T) {
	var err error
	t.Helper()

	// run up migrations for creating tables
	m, err := migrate.New("file://../../../../migrations/", dsn)
	if err != nil {
		t.Fatal(err)
	}

	err = m.Up()
	if err != nil {
		t.Fatal(err)
	}

	// populate tables
	execSqlScript(t, db, "../../../../db/actions/postgres/mock/mock.sql")

	// register a cleanup function for when the test is completed
	t.Cleanup(func() {
		// reset all database changes and close the connection
		err = m.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})
}

// execSqlScript is a helper function to execute SQL commands in the file at the given scriptPath.
func execSqlScript(t *testing.T, db *sql.DB, scriptPath string) {
	t.Helper()

	script, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
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
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
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
