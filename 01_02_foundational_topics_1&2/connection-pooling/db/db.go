package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func New() *sql.DB {
	username := "root"
	password := "rootpassword"
	host := "localhost"
	port := "3306"
	dbname := "kv_store"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		username, password, host, port, dbname)

	// Connect to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error in db connection: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error in db ping:", err)
	}

	log.Println("Connected to database successfully!")

	return db
}
