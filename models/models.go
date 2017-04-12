package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lempiy/pizza-app-pq/utils/utils"
	_ "github.com/lib/pq"
)

// Database type adds handlers for standart sql.DB methods
type Database struct {
	db *sql.DB
}

var database Database
var err error

func init() {
	dbinfo := os.Getenv("DATABASE_URL")
	if dbinfo == "" {
		dbinfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			utils.Keys.UsernameDB, utils.Keys.PasswordDB, utils.Keys.NameDB)
	}
	database.db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
}

func (db Database) begin() (transaction *sql.Tx) {
	transaction, err := db.db.Begin()
	if err != nil {
		log.Println(err)
		return nil
	}
	return transaction
}

func (db Database) prepare(query string) (statement *sql.Stmt) {
	statement, err := db.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return nil
	}
	return statement
}

func (db Database) query(query string, arguments ...interface{}) (rows *sql.Rows) {
	rows, err := db.db.Query(query, arguments...)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rows
}

func (db Database) queryRow(query string, arguments ...interface{}) (row *sql.Row) {
	row = db.db.QueryRow(query, arguments...)
	return row
}

//singleQuery miltiple query isolation
func singleQuery(sql string, args ...interface{}) error {
	SQL := database.prepare(sql)
	tx := database.begin()
	_, err := tx.Stmt(SQL).Exec(args...)

	if err != nil {
		log.Println("singleQuery: ", err)
		tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("singleQuery successful")
	}
	return err
}

//singleQuery miltiple query isolation
func insertWithReturningID(sql string, args ...interface{}) (int, error) {
	var lastID int64
	sql = sql + " RETURNING id;"

	row := database.queryRow(sql, args...)
	row.Scan(&lastID)

	id := int(lastID)
	log.Printf("insertWithReturningID: %d\n", id)
	if err != nil {
		log.Println("insertWithReturningID: ", err)
	}
	return id, err
}

//singleQueryWithAffected miltiple query isolation returns affected rows
func singleQueryWithAffected(sql string, args ...interface{}) (int, error) {
	SQL := database.prepare(sql)
	tx := database.begin()
	result, err := tx.Stmt(SQL).Exec(args...)

	affectedCount, err := result.RowsAffected()
	id := int(affectedCount)
	if err != nil {
		log.Println("singleQuery: ", err)
		tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return 0, err
		}
		log.Println("singleQuery successful")
	}
	return id, err
}

//Close func closes DB connection
func Close() {
	database.db.Close()
}
