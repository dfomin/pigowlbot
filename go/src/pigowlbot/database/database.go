package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "test"
	DB_PASSWORD = "test"
	DB_NAME     = "test"
)

// Controller represents controller for database
type Controller struct {
	DataBase *sql.DB
}

// InitDatabase represents database initialization
func Init() *Controller {
        dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
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

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)        
	}
}
