package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/playmixer/corvid/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	log = logger.New("database")
	DB  *sql.DB
)

func Init() {
	var err error
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	tz := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", host, user, password, name, port, tz)
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.ERROR(err.Error())
		panic(err)
	}

	conn, _ := Connect()
	conn.AutoMigrate(&Ping{})
	conn.AutoMigrate(&Device{})
}

func Connect() (*gorm.DB, error) {

	return gorm.Open(postgres.New(postgres.Config{
		Conn: DB,
	}), &gorm.Config{})
}
