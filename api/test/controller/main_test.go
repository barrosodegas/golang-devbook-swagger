package controller

import (
	"api/src/config"
	"api/src/database"
	"api/src/model"
	"api/src/router"
	"api/test/scripts"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// r Load API routes
var r = router.Generate()
var userToken string

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

	userToken = getUserToken()

	code := m.Run()

	scripts.Run(db)

	os.Exit(code)
}

func getUserToken() string {

	body := []byte(`{"email": "user1@gmail.com", "password": "alb1234"}`)

	req, error := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if error != nil {
		log.Fatal(error)
	}
	response := executeRequest(req)

	var dataAuthentication model.DataAuthentication
	if error := json.Unmarshal(response.Body.Bytes(), &dataAuthentication); error != nil {
		log.Fatal(error)
	}
	return dataAuthentication.Token
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
