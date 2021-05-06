// engine.go

package main

import (
	"database/sql"
	"github.com/gorilla/mux"
    "log"
    "net/http"
    "encoding/json"
    "strconv"
    "fmt"
)

//This is the API core
//Some of this methods could migrate to a json / http util layer
//to make them reusable to other API since the logic would be common

//This method initializes the database connection
//and the routing process to make our api visible to external requests
func (a *App) Initialize(user, password, dbname string) {

	//Sqlite3 avoided, does have risk of file corruption
	//Better performance against postgres although needs consistency check everytime that we start our api

	//sqlite3Conn, _ := sql.Open("sqlite3", "./sqlite-database.db")
	//a.DB = sqlite3Conn
	//a.Router = mux.NewRouter()
    //a.initializeRoutes()

    //Since this is only an exercise, it was used the postgres database to create the necessary tables
    //It should have an unique database if it was a production API
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
    a.initializeRoutes()
}

//This method will allow uas to use frontend dir as our main directory for static / basic frontend layer
//Since we have present a index.html, that will be our main page of frontend
func (a *App) Run(addr string) {

	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/")))

	log.Fatal(http.ListenAndServe(addr, a.Router))
}

//our answers will be JSONFormatted
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//This are the routes defined for our exercise
//They will allow to get statistics about users by/per country
//We will have a csv loader to populate our database user table
//We will have a user creating and search by id or country
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/statistics", a.getNumOfUsersPerCountry).Methods("GET")
	a.Router.HandleFunc("/statistics/{country:[\\w]+}", a.getNumOfUsersByCountry).Methods("GET")
	
	a.Router.HandleFunc("/loader", a.loadDataFromCSV).Methods("GET")
	
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/users", a.getUsersByCountry).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/allUsers", a.getUsers).Methods("GET")

	a.Router.HandleFunc("/countries", a.getCountries).Methods("GET")
}

//This section represents our handlers for our routes
func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Exec getUser")

	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])

	//userId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	row := userJSONFormatted{ID: userId}
	if err := row.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	respondWithJSON(w, http.StatusOK, row)
}

func (a *App) getCountries(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Exec getCountries")

	var countries []customQueryGetCountriesJSONFormatted

	countries, err := getCountries(a.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "No countries found.")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, countries)
}

func (a *App) getUsersByCountry(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Exec getUsersByCountry")

	countryName := r.FormValue("country")
	lines, _ := strconv.Atoi(r.FormValue("lines"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))

	var listOfUsers []userJSONFormatted
	listOfUsers, err := getUsersByCountry(a.DB,countryName,lines,offset)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Users not found for Country "+ countryName)
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, listOfUsers)
}

func (a *App) getNumOfUsersPerCountry(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Exec getNumOfUsersPerCountry")

	var listOfUsersPerCountry []customQueryNumOfUsersPerCountryJSONFormatted

	listOfUsersPerCountry, err := getNumOfUsersPerCountry(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, listOfUsersPerCountry)
}

func (a *App) getNumOfUsersByCountry(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Exec getNumOfUsersByCountry")
	countryName := r.FormValue("country")

	var numOfUsers customQueryNumOfUsersByCountryJSONFormatted
	if err := numOfUsers.getNumOfUsersByCountry(a.DB,countryName); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Users not found for Country "+countryName)
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, numOfUsers)
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Exec getUsers")

	var users []userJSONFormatted

	users, err := getUsers(a.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "No users found.")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Exec getUsers")

	var row userJSONFormatted
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&row); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := row.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, row)
}

func (a *App) loadDataFromCSV(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Exec loadDataFromCSV")

	csvPath := r.FormValue("path")

	var location string = "remote"
	if csvPath == "" {
		location = "local"
	} 

	count, err := loadUsersFromCsv(a,csvPath,location)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, count)
}
