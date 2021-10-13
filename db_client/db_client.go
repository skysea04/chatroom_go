package db_client

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitialiseDBConnection() {
	db, err := sql.Open("mysql", "root:root@tcp(db:3306)/chatroom?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	db.SetConnMaxIdleTime(20)
	db.SetMaxOpenConns(200)

	DB = db
}
