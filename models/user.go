package models

import "log"

//CreateUser creates a new user
func CreateUser(username, password, email string) error {
	sqlQuery := `INSERT INTO person(username, password, email, created_date)
							VALUES($1,$2,$3,now())`
	err := singleQuery(sqlQuery, username, password, email)
	return err
}

//ValidUser check if user data correct
func ValidUser(username, password string) bool {
	var passwordFromDB string
	sqlQuery := `SELECT password
							FROM person WHERE username=$1`
	log.Print("validating user ", username)
	rows := database.query(sqlQuery, username)

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&passwordFromDB)
		if err != nil {
			return false
		}
	}

	if password == passwordFromDB {
		return true
	}

	return false
}

//GetUserID will get the user's ID from the database
func GetUserID(username string, email string) (int, error) {
	var userID int
	sqlQuery := `SELECT id
                        FROM person WHERE username=$1 OR email=$2`
	rows := database.query(sqlQuery, username, email)

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&userID)
		if err != nil {
			return -1, err
		}
	}
	return userID, nil
}

//GetUserIDbyName will get the user's ID from the database only by name
func GetUserIDbyName(username string) (int, error) {
	var userID int
	sqlQuery := `SELECT id
                        FROM person WHERE username=$1`
	rows := database.query(sqlQuery, username)

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&userID)
		if err != nil {
			return -1, err
		}
	}
	return userID, nil
}
