// dataTypes.go

package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	//_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	_ "github.com/lib/pq" //choosed only cause of consistency level against sqlite3
)

//This is core structure that holds database connection and router handler

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//This is the JSON format data type section to
//organize and retrieve tagged fields for better 
//identification on frontend side

type userJSONFormatted struct {
	ID int `json:"userId"`
	Email string `json:"userEmail"`
	Phone string `json:"userPhone"`
	ParcelWeight float64 `json:"userParcelWeight"`
	Country int `json:"countryId"`
}

type customQueryNumOfUsersPerCountryJSONFormatted struct {
	Amount int `json:"count"`
	Country string `json:"countryName"`
}

type customQueryNumOfUsersByCountryJSONFormatted struct {
	Amount int `json:"count"`
	Country string `json:"countryName"`
}

type customQueryGetCountriesJSONFormatted struct {
	Country string `json:"countryName"`
}

//Data type section  

type userDataType struct {
	ID int 
	Email string
	Phone string
	ParcelWeight float64
	Country int
}