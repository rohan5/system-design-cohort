package store

import (
	"errors"
	"log"

	"kv_store/stage_0/internal/db"
)

// KVStore represents a simple in-memory key-value store.
func Put(key string, value any, ttl *int) int {
	db := db.New()
	result, err := db.Exec("INSERT INTO kv_store (k, value) VALUES (?,?) ON DUPLICATE KEY UPDATE value = VALUES(value)", key, value)
	if err != nil {
		log.Fatal("Failed to insert key-value pair: ", err)
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Fatal("Failed to get affected rows: ", err)
	}

	return int(affectedRows)
}

func Get(key string) (string, error) {
	return "", errors.New("not implemented")
}

func Delete(key string) error {
	return errors.New("not implemented")
}
