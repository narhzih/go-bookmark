package e2e

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/pkg/api"
	"gitlab.com/trencetech/mypipe-api/pkg/service"
	"gitlab.com/trencetech/mypipe-api/pkg/service/mailer"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	handler  *http.Server
	logger   = zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	appPort  = "5555"
	pgPort   = 5432
	pgDbName = "mypipe_db"
	pgDbUser = "narhzih"
	pgDbPass = "password"
)

func TestMain(main *testing.M) {
	dbConn, _ := makeDBConn()
	handler = makeHandler(dbConn)
	code := main.Run()
	os.Exit(code)
}

func TestDefaultRoute(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/v1/test-route", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
	//fmt.Println("All well and good")
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.Handler.ServeHTTP(rr, req)

	return rr
}

func makeHandler(dbs db.Database) *http.Server {
	var jwtConfig service.JWTConfig
	jwtConfig.ExpiresIn = 15
	jwtConfig.Key = "hello_world"
	jwtConfig.Algo = jwt.SigningMethodHS256
	//var err error
	apiService := service.NewService(dbs, jwtConfig, initMailer(logger))
	apiHandler := api.NewHandler(apiService, logger)
	router := gin.New()
	rg := router.Group("/v1")
	apiHandler.Register(rg)
	addr := fmt.Sprintf(":%s", appPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	return srv
}

func makeDBConn() (db.Database, error) {
	dbConfig := db.Config{
		Host:           "localhost",
		Port:           pgPort,
		DbName:         pgDbName,
		Username:       pgDbUser,
		Password:       pgDbPass,
		ConnectionMode: "disable",
		Logger:         logger,
	}
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.DbName, dbConfig.ConnectionMode)
	return db.Connect(connectionString, logger)
}

func initMailer(logger zerolog.Logger) *mailer.Mailer {
	var mailerP *mailer.Mailer
	var config mailer.MailConfig
	config.Password = os.Getenv("MAIL_PASSWORD")
	config.Username = os.Getenv("MAIL_USERNAME")
	config.SmtpHost = os.Getenv("MAIL_HOST")
	config.SmtpPort = os.Getenv("MAIL_PORT")
	config.MailFrom = "My Pipe Desk <desk@mypipe.app>"
	logger.Info().Msg(config.Username + " is the username")
	logger.Info().Msg(config.Password + " is the password")
	mailerP = mailer.NewMailer(config)
	return mailerP
}

func applyMigration() {}
