package mydb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect() (*sql.DB, error) {
	// Function to connect to a database
	prod := os.Getenv("PRODUCTION")
	var err error
	if prod != "" {
		mysqlCfg := mysql.Config{
			User:                 os.Getenv("USER"),
			Passwd:               os.Getenv("PASSWORD"),
			Net:                  os.Getenv("NET"),
			Addr:                 os.Getenv("ADDRESS"),
			DBName:               os.Getenv("DB_NAME"),
			AllowNativePasswords: true,
		}
		fmt.Println(mysqlCfg.FormatDSN())
		db, err = sql.Open("mysql", mysqlCfg.FormatDSN())
		if err != nil {
			log.Fatal("Failed to connect", err)
		}
	} else {
		db, err = sql.Open("mysql", os.Getenv("DATABASE"))
		if err != nil {
			log.Fatal(err)
		}
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println("failed on ping", pingErr)
		return nil, fmt.Errorf(pingErr.Error())
	}

	fmt.Println("Connected to db!")
	return db, nil
}
