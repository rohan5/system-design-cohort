package main

import (
	"fmt"
	"kv_store/internal/db"
)

func main() {
	fmt.Println("Creating KV store...")
	db := db.New()

	err := db.Ping()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	// Create the kv_store table
	createStoreTableQuery := `CREATE TABLE IF NOT EXISTS kv_store (
		k VARCHAR(255) PRIMARY KEY,
		value TEXT NOT NULL,
		expires_at TIMESTAMP NULL
	)`

	_, err = db.Exec(createStoreTableQuery)
	if err != nil {
		fmt.Println("Error creating KV store table:", err)
		return
	}
	fmt.Println("KV store created successfully.")

}
