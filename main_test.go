// main_test.go

package main_test

import (
	"os"
	"testing"
    "log"
    "net/http"
    "net/http/httptest"
    "strconv"
    "encoding/json"
    "bytes"

	"github.com/pferreiratlantic/Golang_api"
)

var a main.App

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS country
(
    countryId NUMERIC(4) PRIMARY KEY NOT NULL DEFAULT 0,
    countryName VARCHAR(20) NOT NULL
);
CREATE TABLE IF NOT EXISTS parcelUser
(
    userId INTEGER PRIMARY KEY,
    userEmail TEXT NOT NULL,
    userPhone VARCHAR(16) NOT NULL,
    userParcelWeight NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    countryId NUMERIC(4) NOT NULL DEFAULT 0,
    CONSTRAINT fk_user FOREIGN KEY(countryId) REFERENCES country(countryId)
);`

const tableDefaultEntriesQuery = `
INSERT INTO country(countryId, countryName) VALUES(0,'Unidentified') ON CONFLICT (countryId) DO NOTHING;;
INSERT INTO country(countryId, countryName) VALUES(1,'Cameroon') ON CONFLICT (countryId) DO NOTHING;;
INSERT INTO country(countryId, countryName) VALUES(2,'Ethiopia') ON CONFLICT (countryId) DO NOTHING;;
INSERT INTO country(countryId, countryName) VALUES(3,'Morocco') ON CONFLICT (countryId) DO NOTHING;;
INSERT INTO country(countryId, countryName) VALUES(4,'Mozambique') ON CONFLICT (countryId) DO NOTHING;;
INSERT INTO country(countryId, countryName) VALUES(5,'Uganda') ON CONFLICT (countryId) DO NOTHING;;`

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	ensureTableExists()
	ensureDefaultEntriesExists()

	clearTable("parcelUser")

	code := m.Run()

	os.Exit(code)
}

//DEPRECATED
func ensureDatabaseExists() {
	//file, err := os.Create("sqlite-database.db") // Create SQLite file
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//file.Close()
	//sqlite3Conn, _ := sql.Open("sqlite3", "./sqlite-database.db")
	//a.DB = sqlite3Conn
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableDefaultEntriesQuery); err != nil {
		log.Fatal(err)
	}
}

func ensureDefaultEntriesExists() {
	if _, err := a.DB.Exec(tableDefaultEntriesQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable(tableName string) {
	var query string = "DELETE FROM "+tableName
	if _, err := a.DB.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func clearSequence(tableName string, field string) {
	var query string = "ALTER SEQUENCE "+tableName+"_"+field+"_seq RESTART WITH 1"
	if _, err := a.DB.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

//TEST SECTION TO USER TABLE
//TEST CREATE
//TEST GET
//TEST FAILED GET
//TEST EMPTY USER TABLE
func TestCreateUser(t *testing.T) {

	clearTable("parcelUser")

    var jsonStr = []byte(`{"userEmail":"test@test.pt", 
    						"userPhone": "123456789012", 
    						"userParcelWeight": 1.11,
    						"countryId": 0}`)

    req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["userEmail"] != "test@test.pt" {
		t.Errorf("Expected user email to be 'test@test.pt'. Got '%v'", m["userEmail"])
	}

	if m["userPhone"] != "123456789012" {
		t.Errorf("Expected userPhone to be '123456789012'. Got '%v'", m["userPhone"])
	}

	if m["userParcelWeight"] != 1.11 {
		t.Errorf("Expected userParcelWeight to be '1.11'. Got '%v'", m["userParcelWeight"])
	}
}

func TestGetUser(t *testing.T) {

	clearTable("parcelUser")

	addDummyUser(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	clearTable("parcelUser")
}

func addDummyUser(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO parcelUser(userID, userEmail, userPhone, userParcelWeight, countryId) VALUES($1, $2, $3, $4, $5);", 1, "test"+strconv.Itoa(i)+"@test.pt", "123123123123" , (i+1.0)*10, 0)
	}
}

func TestGetNonExistentUser(t *testing.T) {
	clearTable("parcelUser")

	req, _ := http.NewRequest("GET", "/user/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "User not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'user not found'. Got '%s'", m["error"])
	}
}

func TestEmptyUserTable(t *testing.T) {
	clearTable("parcelUser")

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

//TEST SECTION TO COUNTRY TABLE
//TEST GET

func TestGetCountries(t *testing.T) {

	req, _ := http.NewRequest("GET", "/countries", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestEmptyCountryTable(t *testing.T) {
	clearTable("country")

	req, _ := http.NewRequest("GET", "/countries", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}

	ensureDefaultEntriesExists()
}