package main;

import (
    "database/sql"
    "fmt"
    "log"
    "strings"

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

    // 
    var values []string
    for i := 1; i <= 9; i++ {
        values = append(values, fmt.Sprintf("('%d_a')", i))
        values = append(values, fmt.Sprintf("('%d_b')", i))
        values = append(values, fmt.Sprintf("('%d_c')", i))
        values = append(values, fmt.Sprintf("('%d_d')", i))
        values = append(values, fmt.Sprintf("('%d_e')", i))
        values = append(values, fmt.Sprintf("('%d_f')", i))
    }

    seed_seats_query := "INSERT INTO seats (seat) values " + 
        strings.Join(values, ",")

    result, err := db.Exec(seed_seats_query);

    if err != nil {
        log.Fatal("Error in inserting into seats table ", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Fatal("Error in fetching rows affected ", err)
    } 
    
    log.Printf("Seats table seeded successfully! Rows affected: %d", rowsAffected)

}