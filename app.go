// app.go

package main

import (
	"database/sql"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// tom: initial function is empty, it's filled afterwards
// func (a *App) Initialize(user, password, dbname string) { }

// tom: added "sslmode=disable" to connection string
func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

    // tom: this line is added after initializeRoutes is created later on
    a.initializeRoutes()
}

// tom: initial version
// func (a *App) Run(addr string) { }
// improved version
func (a *App) Run(addr string) {

	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// tom: these are added later
func (a *App) getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	row := customer{ID: id}
	if err := row.getCustomer(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Customer not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, row)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func (a *App) getCustomers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	customers, err := getCustomers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, customers)
}

func (a *App) createCustomer(w http.ResponseWriter, r *http.Request) {
	var row customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&row); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := row.createCustomer(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, row)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/customers", a.getCustomers).Methods("GET")
	a.Router.HandleFunc("/customer", a.createCustomer).Methods("POST")
	a.Router.HandleFunc("/customer/{id:[0-9]+}", a.getCustomer).Methods("GET")
}
