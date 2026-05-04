package main;

import (
	"fmt"
	"log"
	"time"
	"math/rand"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 0. Make db connection
	// 1. take user input for user id
	// 2. take seat number input
	// 3. book seat that user says
	// 4. fetch seats from db and log in nice format

	// make mysql connection
	db := getDbConnection()
    defer db.Close()
	
	// take user input for user id
	var userId int
	fmt.Println("Hi, Welcome! Enter your user id to continue.")
	fmt.Scanln(&userId)
	log.Printf("Hello, %d! Let's get started.\n", userId)

	// assign random seat instead of taking input from user
	// first fetch all seats from db and select 1 from them randomly and assign to user
	var seats = getAllSeats(db)
	rand.Seed(time.Now().UnixNano())
	seat := seats[rand.Intn(len(seats))]

	logSeatsStatus(db)

	// book seat that user says
	bookSeat(db, seat, userId)

	logSeatsStatus(db)
	

}

func getAllSeats(db *sql.DB) []string {
	rows, err := db.Query("SELECT * FROM seats")
	if err != nil {
		log.Fatal("Error querying seats:", err)
	}
	defer rows.Close()
	var seats []string
	for rows.Next(){
		var id int
		var seat string
		var customer_id sql.NullInt64
		err := rows.Scan(&id, &seat, &customer_id)
		if err != nil {
			log.Fatal("Error scanning seat row:", err)
		}
		seats = append(seats, seat)
	}
	return seats
}

func bookSeat(db *sql.DB, seat string, userId int) error {
	updateQuery := "UPDATE seats SET customer_id = ? WHERE seat = ?;"
	result, err := db.Exec(updateQuery, userId, seat)
	if err != nil {
		log.Printf("Error booking seat %s for user %d: %v\n", seat, userId, err)
	} else {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Error fetching rows affected for seat %s booking by user %d: %v\n", seat, userId, err)
		}
		log.Println("Seat booking resulst! >> Rows affected: ", rowsAffected)
	}
	return err
}

func logSeatsStatus(db *sql.DB) {
// log seat booking result
	rows, err := db.Query("SELECT * FROM seats")
	if err != nil {
		log.Fatal("Error querying seats:", err)
	}
	defer rows.Close()
	// Then iterate with rows.Next(), rows.Scan(), etc.
	// For simplicity, we will just log the seat numbers here.
	log.Println("Seats status:")

	// instead of just logging, 
	// format it in a way similar to shown in flight booking  apps
	// in a 2D format with seat numbers and user ids (if booked)
	// example format:
	// x  21  x  x
	// x  x  x  x
	// x  x  34  x
	// where x is for empty seats and user_id if seat iss taken.
	i := 1
	// j := 1
	// var seatMap [][]string
	for rows.Next() {
		var id int
		var seat string
		var bookedBy sql.NullInt64 // Use NullInt64 to handle NULL values for user_id
		err := rows.Scan(&id, &seat, &bookedBy)
		if err != nil {
			log.Fatal("Error scanning seat row:", err)
		}
		if bookedBy.Valid {
			fmt.Print(bookedBy.Int64, " ")
		} else {
			fmt.Print("x", " ")
		}
		// if i%6 == 0 {
		// 	fmt.Println() // New line after every 6 seats for better formatting
		// }
		if i%6 == 0 {
			fmt.Println() // New space after every 3 seats for better formatting
		} else if i%3 == 0 {
			fmt.Print(" ") // New space after every 3 seats for better formatting
		}
		if i%18 == 0 {
			fmt.Println() // New line after every 5 lines of seats for better formatting
		}
		
		i++
		// fmt.Printf("Seat ID: %d, Seat Number: %s, Booked By: %v\n", id, seat, bookedBy.Int64)

		// store in 2d array to log later in formatted way
		// value_to_store := "x"
		// if bookedBy.Valid {
		// 	value_to_store = fmt.Sprintf("%d", bookedBy.Int64)
		// }

	}
}

func getDbConnection() *sql.DB {
	username := "root"
    password := "rootpassword"
    host := "localhost"
    port := "3306"
    dbname := "airline_booking"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		username, password, host, port, dbname)

	 // Connect to the database
    db,err := sql.Open("mysql", dsn);
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