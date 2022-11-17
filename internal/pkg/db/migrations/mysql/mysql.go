package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

var Db *sql.DB

func InitDB() {
	db, err := sql.Open("mysql", "root:mysql@tcp(localhost)/hackernews")
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}

	Db = db
}

func CloseDB() error {
	return Db.Close()
}

func Migrate() {
	if err := Db.Ping(); err != nil {
		log.Fatalln(err)
	}

	driver, err := mysql.WithInstance(Db, &mysql.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/mysql",
		"hackernews",
		driver,
	)
	if err != nil {
		log.Fatalln(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalln(err)
	}
}
