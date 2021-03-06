package controller

import (
	"api/src/config"
	"api/src/database"
	"api/src/model"
	"api/src/router"
	"api/test/scripts"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// r Load API routes
var r = router.Generate()
var userToken string
var userTokenToUpdate string

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

	userToken = GetUserToken("user1@gmail.com", "alb1234")
	userTokenToUpdate = GetUserToken("user2@gmail.com", "alb1234")

	code := m.Run()

	scripts.Run(db)

	os.Exit(code)
}

func GetUserToken(email, password string) string {

	body := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password))

	response := executeRequest("POST", "/login", "", body, true)

	var dataAuthentication model.DataAuthentication
	if error := json.Unmarshal(response.Body.Bytes(), &dataAuthentication); error != nil {
		log.Fatal(error)
	}
	return dataAuthentication.Token
}

// executeRequest Executes the request to be tested and returns a response.
func executeRequest(method, uri, token string, body []byte, hasContentTypeJson bool) *httptest.ResponseRecorder {

	var requestBody io.Reader = nil
	if body != nil {
		requestBody = bytes.NewBuffer(body)
	}

	req, error := http.NewRequest(method, uri, requestBody)
	if error != nil {
		log.Fatal(error)
	}

	if hasContentTypeJson {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

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
