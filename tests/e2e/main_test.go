package e2e

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	handler *http.Server
)

func TestMain(main *testing.M) {
	code := main.Run()
	os.Exit(code)
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

func compareResponse() {}
