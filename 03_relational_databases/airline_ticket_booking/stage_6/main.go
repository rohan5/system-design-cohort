package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

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

	PreCleanup(db)
	startTime := time.Now()

	// assign random seat instead of taking input from user
	// first fetch all seats from db and select 1 from them randomly and assign to user
	// var seats = getAllSeats(db)

	// fetch all customers and book random seat in parallel for each customer
	rand.Seed(time.Now().UnixNano())
	var customers = getAllCustomers(db)
	var wg sync.WaitGroup
	wg.Add(len(customers))
	for _, customerId := range customers {
		go func(customerId int) {
			defer wg.Done()
			//seat := seats[rand.Intn(len(seats))]

			// log.Printf("Trying to book seat %s for user %d...\n", seat, customerId)
			res, err := bookSeat(db, customerId)
			if err != nil {
				log.Printf("Error booking seat for user %d: %v\n", customerId, err)
			} else {
				log.Printf("Successfully booked seat %s for user %d\n", *res, customerId)
			}
		}(customerId)
	}
	wg.Wait()

	log.Printf("Time taken: %s\n", time.Since(startTime))
	// logSeatsStatus(db)

	// book seat that user says

	logSeatsStatus(db)

}

func getAllCustomers(db *sql.DB) []int {
	var customers []int
	rows, err := db.Query("SELECT * FROM customers ORDER BY id ASC")
	if err != nil {
		log.Fatal("Error querying customers:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal("Error scanning customer row:", err)
		}
		customers = append(customers, id)

	}
	return customers
}

func getAllSeats(db *sql.DB) []string {
	rows, err := db.Query("SELECT * FROM seats")
	if err != nil {
		log.Fatal("Error querying seats:", err)
	}
	defer rows.Close()
	var seats []string
	for rows.Next() {
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

func bookSeat(db *sql.DB, userId int) (*string, error) {
	trxn, _ := db.Begin()

	row := trxn.QueryRow(`SELECT id, seat, customer_id FROM seats WHERE customer_id IS null ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED`)
	if row.Err() != nil {
		log.Printf("Error fetching available seat for user %d: %v\n", userId, row.Err())
		return nil, row.Err()
	}

	var seatId int
	var seat string
	var bookedBy sql.NullInt64

	err := row.Scan(&seatId, &seat, &bookedBy)
	if err != nil {
		trxn.Rollback()
		return nil, err
	}

	result2, err := trxn.Exec(`UPDATE seats SET customer_id = ? WHERE id = ?`, userId, seatId)
	if err != nil {
		return nil, err
	}

	err = trxn.Commit()
	if err != nil {
		return nil, err
	}
	log.Printf("Rows affected: %d", result2)

	return &seat, nil
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

func PreCleanup(db *sql.DB) {
	db.Exec(`UPDATE seats SET customer_id = null`)
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
