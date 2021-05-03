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

)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM customer")
	a.DB.Exec("ALTER SEQUENCE customer_customerId_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS customer
(
    customerId SERIAL,
    customerEmail TEXT NOT NULL,
    customerPhone VARCHAR(12) NOT NULL,
    customerParcelWeight NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT customer_pkey PRIMARY KEY (customerId)
);`

// tom: next functions added later, these require more modules: net/http net/http/httptest
func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/customers", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
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

func TestGetNonExistentCustomer(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/customer/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Customer not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Customer not found'. Got '%s'", m["error"])
	}
}

// tom: rewritten function
func TestCreateCustomer(t *testing.T) {

	clearTable()

    var jsonStr = []byte(`{"customerEmail":"test@test.pt", 
    						"customerPhone": "123456789012", 
    						"customerParcelWeight": 1.11}`)

    req, _ := http.NewRequest("POST", "/customer", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["customerEmail"] != "test@test.pt" {
		t.Errorf("Expected customer email to be 'test@test.pt'. Got '%v'", m["customerEmail"])
	}

	if m["customerPhone"] != "123456789012" {
		t.Errorf("Expected customerPhone to be '123456789012'. Got '%v'", m["customerPhone"])
	}

	if m["customerParcelWeight"] != 1.11 {
		t.Errorf("Expected customerParcelWeight to be '1.11'. Got '%v'", m["customerParcelWeight"])
	}
}


func TestGetCustomer(t *testing.T) {
	clearTable()
	addCustomer(10)

	req, _ := http.NewRequest("GET", "/customer/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addCustomer(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO customer(customerEmail, customerPhone, customerParcelWeight) VALUES($1, $2, $3);", "test"+strconv.Itoa(i)+"@test.pt", "123123123123" , (i+1.0)*10)
		
	}
}