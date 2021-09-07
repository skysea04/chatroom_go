package db_client

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DBClient *sql.DB

func InitialiseDBConnection() {
	db, err := sql.Open("mysql", "root:root@/chatroom")
	if err != nil {
		panic(err.Error())
	}
	db.SetConnMaxIdleTime(20)
	db.SetMaxOpenConns(200)
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	DBClient = db
}
