package e2e

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/mypipeapp/mypipeapi/cmd/api/internal"
	"github.com/mypipeapp/mypipeapi/cmd/api/routes"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var (
	handler           *http.Server
	db                *sql.DB
	dsn               string
	app               internal.Application
	logger            = zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	globalAccessToken string
	globalUserID      int
	globalUserEmail   string
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		logger.Err(err).Msg("Could not load environment variables")
		os.Exit(1)
	}

	// create application instance and store
	// it globally
	dsn = os.Getenv("DB_DSN_TEST")
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// run up migrations for creating tables
	mig, err := migrate.New("file://../../../../migrations/", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = mig.Up()
	if err != nil {
		log.Fatal(err)
	}
	// populate tables
	execSqlScript(db, "../../../../db/actions/postgres/mock/mock.sql")

	app = createApplicationInstance()

	// setup router
	router := gin.Default()
	rg := router.Group("/v1")
	routes.BootRoutes(app, rg)

	appPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("unable to bind port")
	}
	addr := fmt.Sprintf(":%d", appPort)

	handler = &http.Server{
		Addr:    addr,
		Handler: router,
	}
	createGlobalUserAndLogin()
	code := m.Run()
	err = mig.Drop()
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
	os.Exit(code)
}

func TestHealthZ(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/healthz", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
}
