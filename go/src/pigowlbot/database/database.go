package database

import (
	"database/sql"
	"fmt"
	"log"
	"pigowlbot/private"
	_ "github.com/lib/pq"
)

// Controller represents controller for database
type Controller struct {
	DataBase *sql.DB
}

// InitDatabase represents database initialization
func Init() *Controller {
        dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", private.DB_USER, private.DB_PASSWORD, private.DB_NAME)
        db, err := sql.Open("postgres", dbinfo)
        checkErr(err)

	err = db.Ping()
	checkErr(err)

	return &Controller{DataBase: db}
}

func (c *Controller) AddSubscriber(chatId int64) {
	stmt, err := c.DataBase.Prepare("INSERT INTO chat(chat_id) VALUES($1)")
	checkErr(err)
	_, err = stmt.Exec(chatId)
	checkErr(err)
}

func (c *Controller) CheckSubscriber(chatId int64) bool {
	var exists bool
	query := "SELECT exists (SELECT * FROM chat where chat_id=$1)"
	err := c.DataBase.QueryRow(query, chatId).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("error checking if row exists '%s' %v", chatId, err)
	}
	return exists
}

func (c *Controller) GetSubscribers() []int64 {
	rows, err := c.DataBase.Query("SELECT * FROM golang.chat")
	checkErr(err)

	var result []int64
	for rows.Next() {
		var id int
		var chatId int64
		err = rows.Scan(&id, &chatId)
		checkErr(err)
		result = append(result, chatId)
	}
	return result
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
