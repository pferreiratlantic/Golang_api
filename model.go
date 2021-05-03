// model.go
package main

import (
	"database/sql"
)


// tom: add backticks to json
type customer struct {
	ID    int     `json:"customerId"`
	Email  string  `json:"customerEmail"`
	Phone string `json:"customerPhone"`
	ParcelWeight float64 `json:"customerParcelWeight"`
}

// tom: these are added after tdd tests
func (row *customer) getCustomer(db *sql.DB) error {
	return db.QueryRow("SELECT customerEmail, customerPhone, customerParcelWeight FROM customer WHERE customerId=$1",
		row.ID).Scan(&row.Email, &row.Phone, &row.ParcelWeight)
}

func getCustomers(db *sql.DB, start, count int) ([]customer, error) {
	rows, err := db.Query(
		"SELECT customerId, customerEmail, customerPhone, customerParcelWeight FROM customer LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	customers := []customer{}

	for rows.Next() {
		var row customer
		if err := rows.Scan(&row.ID, &row.Email, &row.Phone, &row.ParcelWeight); err != nil {
			return nil, err
		}
		customers = append(customers, row)
	}

	return customers, nil
}

func (row *customer) createCustomer(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO customer(customerEmail, customerPhone, customerParcelWeight) VALUES($1, $2, $3) RETURNING customerId",
		row.Email, row.Phone, row.ParcelWeight).Scan(&row.ID)

	if err != nil {
		return err
	}

	return nil
}