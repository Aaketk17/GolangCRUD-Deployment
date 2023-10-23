package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func PGConnection() *pgx.Conn {
	envErr := godotenv.Load()

	if envErr != nil {
		fmt.Println("Error in loading .env file, ", envErr)
	}

	host := os.Getenv("DB_URL")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	portInt, er := strconv.Atoi(port)
	if er != nil {
		fmt.Println("Error in converting string to number", er)
	}

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	pgConnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres", host, portInt, user, password)
	fmt.Println(pgConnString)

	pgConn, err := pgx.Connect(ctx, pgConnString)
	if err != nil {
		fmt.Printf("Unable to connection to database: \n%v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Database connection establised to db ", dbname)
	}

	return pgConn
}

var Connection *pgx.Conn = PGConnection()
