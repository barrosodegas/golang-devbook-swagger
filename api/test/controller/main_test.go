package controller

import (
	"api/src/config"
	"api/src/database"
	"api/src/router"
	"api/test/scripts"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// r Load API routes
var r = router.Generate()

// TestMain load environment variables, database and run package tests
func TestMain(m *testing.M) {

	if error := os.Setenv("ENV_PATH", "../.env-test"); error != nil {
		log.Fatal(error)
	}

	config.Load()

	db, error := database.Connect()
	if error != nil {
		log.Fatal("Connection failed.")
		return
	}

	defer db.Close()

	scripts.Run(db)

	code := m.Run()

	scripts.Run(db)

	os.Exit(code)
}

// executeRequest Executes the request to be tested and returns a response.
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode Checks that the endpoint responded as expected.
func checkResponseCode(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
