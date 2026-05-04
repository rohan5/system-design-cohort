package main;

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

func main() {

    username := "root"
    password := "rootpassword"
    host := "localhost"
    port := "3306"
    dbname := "airline_booking"


    // MySQL connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		username, password, host, port, dbname)


    // Connect to the database
    db,err := sql.Open("mysql", dsn);
    if err != nil {
        log.Fatal("Error in db connection: ", err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal("Error in db ping:", err)
    }

    log.Println("Connected to database successfully!")

    create_seats_table_query := "CREATE TABLE IF NOT EXISTS seats (" +
        "id int AUTO_INCREMENT PRIMARY KEY," +
        "seat varchar(255) NOT NULL," +
        "customer_id int" +
        ");"

    result, err := db.Exec(create_seats_table_query);

    if err != nil {
        log.Fatal("Error in creating seats table ", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Fatal("Error in fetching rows affected ", err)
    }
    
    log.Printf("Seats table created successfully! Rows affected: %d", rowsAffected)

}