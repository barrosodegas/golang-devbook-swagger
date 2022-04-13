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

var r = router.Generate()

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

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestLoginWithSucces(t *testing.T) {

	body := []byte(`{"email": "user1@gmail.com", "password": "alb1234"}`)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var dataAuthentication model.DataAuthentication
	if error := json.Unmarshal(response.Body.Bytes(), &dataAuthentication); error != nil {
		t.Errorf("Expected an DataAuthentication model. Got %s", body)
	}
}

func TestLoginWithBadRequestError(t *testing.T) {

	body := []byte(``)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestLoginWithUnauthorizedError(t *testing.T) {

	body := []byte(`{"email": "user1@gmail.com", "password": "alb12"}`)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}