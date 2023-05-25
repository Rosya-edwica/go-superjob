package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"superjob/pkg/logger"
	"superjob/pkg/telegram"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DB struct {
	db   *sql.DB
	Name string
	Host string
	Port int
	User string
	Pass string
}

func InitDatabase() *DB {
	err := godotenv.Load()
	checkErr(err)
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	checkErr(err)
	db := &DB{
		Host: os.Getenv("POSTGRES_HOST"),
		User: os.Getenv("POSTGRES_USER"),
		Name: os.Getenv("POSTGRES_DATABASE"),
		Pass: os.Getenv("POSTGRES_PASSWORD"),
		Port: port,
	}
	psqlUrl := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, db.User, db.Pass, db.Name)
	connection, err := sql.Open("postgres", psqlUrl)
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
