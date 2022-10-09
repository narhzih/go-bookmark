package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
	"gitlab.com/trencetech/mypipe-api/cmd/api/handlers"
	"gitlab.com/trencetech/mypipe-api/cmd/api/services"
	"gitlab.com/trencetech/mypipe-api/cmd/api/services/mailer"
	"gitlab.com/trencetech/mypipe-api/db"
	"gitlab.com/trencetech/mypipe-api/pkg/service"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"sync"
	"testing"
)

var (
	handler           *http.Server
	logger            = zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	appPort           = "5555"
	pgPort            = 5432
	pgDbName          = "mypipe_db_test"
	pgDbUser          = "narhzih"
	pgDbPass          = "me.password_"
	dbUrl             = ""
	globalAccessToken = ""
	hostAndPort       = "localhost:5432"
)

func TestMain(main *testing.M) {
	dbConn, _ := makeDBConn()
	dbUrl = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgDbUser, pgDbPass, hostAndPort, pgDbName)
	applyMigration("down")
	applyMigration("up")

	handler = makeHandler(dbConn)
	createInitialUser()
	applyMigration("down")
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
		t.Errorf("Expected response codee %d. Got %d\n", expected, actual)
	}
}
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.Handler.ServeHTTP(rr, req)

	return rr
}

func makeHandler(dbs db.Database) *http.Server {
	var jwtConfig services.JWTConfig
	jwtConfig.ExpiresIn = 15
	jwtConfig.Key = "hello_world"
	jwtConfig.Algo = jwt.SigningMethodHS256
	//var err error
	apiService := service.NewService(dbs, jwtConfig, initMailer(logger))
	apiHandler := handlers.NewHandler(apiService, logger)
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

func applyMigration(direction string) {
	var syncer sync.WaitGroup
	var cmd *exec.Cmd
	cmd = exec.Command("migrate", "-database", dbUrl, "-path", "../../migrations", direction)
	if direction == "down" {
		log.Print("Executing migrate down")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		syncer.Add(1)
		go func() {
			defer syncer.Done()
			defer stdin.Close()
			io.WriteString(stdin, "y")
		}()
		syncer.Wait()
		return
	}

	cmd.Stderr = logger
	cmd.Stdout = logger
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
func createInitialUser() {
	// Create user
	body := []byte(`{"username": "dummy_user", "profile_name": "dummy_user", "password": "test123", "email": "dummy@user.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/v1/sign-up", bytes.NewBuffer(body))
	res := executeRequest(req)
	vResInJson := struct {
		Message string `json:"message"`
		Data    struct {
			VToken string `json:"v_token"`
		}
	}{}
	respBody, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(respBody, &vResInJson); err != nil {
		log.Print("Error coming from parsing vResINJson")
		log.Fatal(err)
	}

	log.Print(fmt.Sprintf("Verification token is -> %+v", vResInJson.Message))

	// Verify account
	reqUrl := fmt.Sprintf("/v1/verify-account/%+v", vResInJson.Data.VToken)
	req, _ = http.NewRequest(http.MethodPost, reqUrl, nil)
	res = executeRequest(req)
	lResInJson := struct {
		Message string `json:"message"`
		Data    struct {
			Token string `json:"token"`
			User  struct {
				ID int64 `json:"id"`
			}
		}
	}{}
	respBody, _ = io.ReadAll(res.Body)
	log.Print(fmt.Sprintf("response body log %+v", res.Body))
	if err := json.Unmarshal(respBody, &lResInJson); err != nil {
		log.Print("Error coming from parsing lResINJsonm")
		log.Fatal(err)
	}
	// Login

	globalAccessToken = lResInJson.Data.Token
	log.Print(lResInJson.Message)
	log.Print(globalAccessToken)
}
