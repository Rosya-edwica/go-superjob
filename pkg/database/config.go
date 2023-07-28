package database

import (
	"database/sql"
	"fmt"
	"os"
	"superjob/pkg/logger"
	"superjob/pkg/telegram"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db   *sql.DB
	Name string
	Host string
	Port string
	User string
	Pass string
}

func InitDatabase() *DB {
	err := godotenv.Load()
	checkErr(err)
	checkErr(err)
	db := &DB{
		Host: os.Getenv("MYSQL_HOST"),
		User: os.Getenv("MYSQL_USER"),
		Name: os.Getenv("MYSQL_DATABASE"),
		Pass: os.Getenv("MYSQL_PASSWORD"),
		Port: os.Getenv("MYSQL_PORT"),
	}
	connection, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.User, db.Pass, db.Host, db.Port, db.Name))
	checkErr(err)

	db.db = connection
	return db
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) GetDB() *sql.DB {
	return d.db
}

func checkErr(err error) {
	if err != nil {
		telegram.Mailing(err.Error())
		logger.Log.Fatalln(err.Error())
	}
}
