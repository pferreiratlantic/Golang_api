// dataModel.go

package main

import (
	"database/sql"
	"fmt"
	"strings"
)

// This section represents the Model Layer with all the definitions
// To execute and / or retrieve information from database

//Method responsable for retrieving information from a specific user
func (userJSONFormatted *userJSONFormatted) getUser(db *sql.DB) error {
	return db.QueryRow("SELECT userEmail, userPhone, userParcelWeight FROM parcelUser WHERE userId=$1",
		userJSONFormatted.ID).Scan(&userJSONFormatted.Email, 
			&userJSONFormatted.Phone, 
			&userJSONFormatted.ParcelWeight)
}

//Method responsable for retrieving countries statistics
func getNumOfUsersPerCountry(db *sql.DB) ([]customQueryNumOfUsersPerCountryJSONFormatted, error) {
	rows, err := db.Query(
		"SELECT Count(*) as number_of_users,countryName from parcelUser inner join country using(countryId) Group By countryName")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	usersPerCountry := []customQueryNumOfUsersPerCountryJSONFormatted{}

	for rows.Next() {
		var rowUsersPerCountry customQueryNumOfUsersPerCountryJSONFormatted
		if err := rows.Scan(&rowUsersPerCountry.Amount, 
								&rowUsersPerCountry.Country); err != nil {
			return nil, err
		}
		usersPerCountry = append(usersPerCountry, rowUsersPerCountry)
	}

	return usersPerCountry, nil

}

//Method responsable for retrieving the amount of users on specific country
func (numOfUsersByCountry *customQueryNumOfUsersByCountryJSONFormatted) getNumOfUsersByCountry(db *sql.DB, countryName string) error {
	return db.QueryRow("SELECT Count(*) as number_of_users,countryName from parcelUser inner join country using(countryId) where countryName = '$1'",
			countryName).Scan(&numOfUsersByCountry.Amount, 
			&numOfUsersByCountry.Country)
}

//Method responsable for retrieving available countries
func getCountries(db *sql.DB) ([]customQueryGetCountriesJSONFormatted,error) {

	rows, err := db.Query("SELECT DISTINCT(countryName) from country")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	countries := []customQueryGetCountriesJSONFormatted{}

	for rows.Next() {
		var rowCountry customQueryGetCountriesJSONFormatted
		if err := rows.Scan(&rowCountry.Country); err != nil {
			return nil, err
		}
		countries = append(countries, rowCountry)
	}

	return countries, nil
}

//Method responsable for retrieving all users from specific country
func getUsersByCountry(db *sql.DB, countryName string,lines int,offset int) ([]userJSONFormatted, error) {
	rows, err := db.Query(
		"SELECT userId, userEmail, userPhone, userParcelWeight, countryId FROM parcelUser inner join country using(countryId) where countryName =$1 ORDER BY userId LIMIT $2 OFFSET $3",
		countryName,lines,offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []userJSONFormatted{}

	for rows.Next() {
		var rowuser userJSONFormatted
		if err := rows.Scan(&rowuser.ID, 
								&rowuser.Email, 
								&rowuser.Phone, 
								&rowuser.ParcelWeight,
								&rowuser.Country); err != nil {
			return nil, err
		}
		users = append(users, rowuser)
	}

	return users, nil
}

//Method responsabli for retrieving all users
func getUsers(db *sql.DB) ([]userJSONFormatted, error) {
	rows, err := db.Query(
		"SELECT userId, userEmail, userPhone, userParcelWeight, countryId FROM parcelUser")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []userJSONFormatted{}

	for rows.Next() {
		var rowuser userJSONFormatted
		if err := rows.Scan(&rowuser.ID, 
								&rowuser.Email, 
								&rowuser.Phone, 
								&rowuser.ParcelWeight,
								&rowuser.Country); err != nil {
			return nil, err
		}
		users = append(users, rowuser)
	}

	return users, nil
}

//Method responsable for creating user
func (rowuser *userJSONFormatted) createUser(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO parcelUser(userID,userEmail, userPhone, userParcelWeight, countryId) VALUES($1, $2, $3, $4, $5) RETURNING userId",
		rowuser.ID,
		rowuser.Email, 
		rowuser.Phone, 
		rowuser.ParcelWeight,
		rowuser.Country).Scan(&rowuser.ID)

	if err != nil {
		return err
	}

	return nil
}

//Method responsable for creating users in bulk
//DEPRECATED :: Low Performance
func createSetOfUsers(db *sql.DB,setOfUsers []userDataType) (int, error) {
	
	var counter int = 0

	for i := 0; i < len(setOfUsers); i++ {
		var rowuser = setOfUsers[i]
		fmt.Printf("Row %d : email: %s phone: %s parcelWeight: %f country: %d \n", i, rowuser.Email,rowuser.Phone,rowuser.ParcelWeight,rowuser.Country)

		err := db.QueryRow(
		"INSERT INTO parcelUser(userID,userEmail, userPhone, userParcelWeight,countryId) VALUES($1, $2, $3, $4, $5)",
		rowuser.ID,
		rowuser.Email, 
		rowuser.Phone, 
		rowuser.ParcelWeight,
		rowuser.Country)

		if err != nil {
			return 0, nil
		}
		counter++
	}

	return counter, nil
}

//Method responsable for creating users in bulk
//Higher Performance : does batches of 8k reg at a time
func createSetOfUsersByBatch(db *sql.DB,setOfUsers []userDataType) (int, error) {
    var numOfColumns int = 5
    var offset int = 8000

    if(len(setOfUsers) < offset){
    	offset = len(setOfUsers)
    }

    var counter int = 0
    valueArgs := make([]interface{}, 0, offset * numOfColumns)
    valueStrings := make([]string, 0, offset)
    i := 0

    for _, post := range setOfUsers {
        valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d , $%d)", i*numOfColumns+1, i*numOfColumns+2, i*numOfColumns+3, i*numOfColumns+4, i*numOfColumns+5))
        valueArgs = append(valueArgs, post.ID)
        valueArgs = append(valueArgs, post.Email)
        valueArgs = append(valueArgs, post.Phone)
        valueArgs = append(valueArgs, post.ParcelWeight)
        valueArgs = append(valueArgs, post.Country)
        counter++
        i++
        if i == offset{

        	stmt := fmt.Sprintf("INSERT INTO parcelUser (userId, userEmail, userPhone, userParcelWeight,countryId) VALUES %s ON CONFLICT (userId) DO NOTHING", 
                        strings.Join(valueStrings, ","))
    		_, err := db.Exec(stmt, valueArgs...)
    		if err != nil{
				return counter, err
			}
        	valueArgs = nil
        	valueStrings = nil
        	i = 0
        }
    }

    stmt := fmt.Sprintf("INSERT INTO parcelUser (userId, userEmail, userPhone, userParcelWeight,countryId) VALUES %s ON CONFLICT (userId) DO NOTHING", 
                strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	if err != nil{
		return counter, err
	}
    
    return len(setOfUsers), nil
}