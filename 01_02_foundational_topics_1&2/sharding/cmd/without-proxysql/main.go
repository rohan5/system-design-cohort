package main

import (
	"database/sql"
	"fmt"
	"sharding/internal/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

var DBs = []*sql.DB{}

func init() {
	db1 := db.New("to_do_list_1")
	DBs = append(DBs, db1)

	db2 := db.New("to_do_list_2")
	DBs = append(DBs, db2)
}

func getDBIndex(userId int) int {
	if userId == 1 {
		return 0
	} else {
		return 1
	}
}

func addTask(userId int, task string) {
	res, err := DBs[getDBIndex(userId)].Exec(`INSERT INTO to_do_list (user_id, task) VALUES (?,?)`, userId, task)
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
