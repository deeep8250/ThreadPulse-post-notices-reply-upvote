package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/threadpulse/internal/config"
)

func DBint() {

	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"

	var err error
	for range 5 {
		config.PostgisDB, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("error while connecting to db : ", err.Error())
	}

	err = config.PostgisDB.Ping()
	if err != nil {
		log.Fatal("database connection failed")

	}

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatal("Migration files not found ", err.Error())
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migtaion failed ", err.Error())
	}
	log.Println("Migrations ran successfully")

	fmt.Println("database connect successfully")

}
