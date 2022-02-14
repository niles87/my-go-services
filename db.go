package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connect() {
	prod := os.Getenv("PRODUCTION")
	var err error
	if prod != "" {
		db, err = sql.Open("mysql", os.Getenv("JAWSDB_URL"))
		if err != nil {
			log.Fatal(err)
		}
	}

	db, err = sql.Open("mysql", os.Getenv("DATABASE"))
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to db!")
}

func queryAllUsers() ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		return nil, fmt.Errorf("queryAllUsers: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Wins, &user.Losses, &user.Draws); err != nil {
			return nil, fmt.Errorf("queryAllUsers: %v", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("queryAllUsers: %v", err)
	}

	return users, nil
}

func addUser(user User) (int64, error) {
	res, err := db.Exec("INSERT INTO user (name, email, password, wins, losses, draws) VALUES (?,?,?,?,?,?)", user.Name, user.Email, user.Password, user.Wins, user.Losses, user.Draws)

	if err != nil {
		return 0, fmt.Errorf("addUser %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addUser %v", err)
	}

	return id, nil
}

func queryUserByID(id int64) (User, error) {
	var user User

	row := db.QueryRow("SELECT * FROM user WHERE id=?", id)

	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Wins, &user.Losses, &user.Draws); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("queryUserById no record with id: %d ", id)
		}
		return user, fmt.Errorf("addUser %v", err)
	}

	return user, nil
}
