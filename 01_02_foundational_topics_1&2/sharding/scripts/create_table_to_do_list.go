package main

import (
	"fmt"
	"sharding/internal/db"
)

func main() {
	fmt.Println("Creating table...")
	db := db.New("to_do_list_1")
	defer db.Close()
	res, err := db.Exec(`CREATE TABLE IF NOT EXISTS to_do_list 
			(id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			task VARCHAR(255) NOT NULL
			);`)

	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		return
	}
	fmt.Println("Table created successfully.", rowsAffected, "rows affected.")
}
