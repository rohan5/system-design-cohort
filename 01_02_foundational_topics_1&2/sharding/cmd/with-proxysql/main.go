package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// var DBs = []*sql.DB{}

func getDBConnection() *sql.DB {
	username := "appuser"
	password := "apppass"
	host := "localhost"
	port := "6033"
	dbname := "app_db"

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

	// log.Println("Connected to database successfully!")

	return db
}

func addShardHint(query string, userId int) string {
	shard := 0
	if userId == 1 {
		shard = 1
	} else {
		shard = 0
	}

	return fmt.Sprintf("/* shard:%d */ %s", shard, query)
}

func addTask(userId int, task string) {
	db := getDBConnection()
	query := addShardHint(`INSERT INTO to_do_list (user_id, task) VALUES (?,?)`, userId)
	res, err := db.Exec(query, userId, task)
	if err != nil {
		fmt.Println("Error adding task:", err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		return
	}
	fmt.Println("Task added successfully.", rowsAffected, "rows affected.")
}

func main() {
	fmt.Println("Hello, World!")

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/add-task/:userId/:task", func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("userId"))
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid userId"})
			return
		}
		task := c.Param("task")

		addTask(userId, task)
		c.JSON(200, gin.H{
			"message": "Task added successfully",
		})
	})

	router.Run()
}
