package postgres

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

const (
	skipMessage = "postgres: skipping integration test"
)

// newTestDb prepares the test database by applying migrations and populating with test data.
// It returns a connection to the test database.
func newTestDb(t *testing.T) *sql.DB {
	t.Helper()

	// establish a connection to the test database
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Fatal(err)
	}

	dsn := os.Getenv("DB_DSN_TEST")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

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
	execSqlScript(t, db, "./testdata/populate_data.sql")

	// register a cleanup function for when the test is completed
	t.Cleanup(func() {
		// reset all database changes and close the connection
		err = m.Drop()
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})
	return db
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
