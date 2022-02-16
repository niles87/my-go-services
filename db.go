package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connect() {
	prod := os.Getenv("PRODUCTION")
	var err error
	if prod != "" {
		mysqlCfg := mysql.Config{
			User:   os.Getenv("USER"),
			Passwd: os.Getenv("PASSWORD"),
			Net:    os.Getenv("NET"),
			Addr:   os.Getenv("ADDRESS"),
			DBName: os.Getenv("DB_NAME"),
		}
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

func deleteUserByID(id int64) (int64, error) {
	stmt, err := db.Prepare("DELETE FROM user WHERE id=?")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %v", err)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("deleteUserByID: %v", err)
	}

	rowsRemoved, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("deleteUserById: %v", err)
	}

	return rowsRemoved, nil
}

func updateUserByID(id int64, user User) (int64, error) {
	stmt, err := db.Prepare("UPDATE user SET name=?, email=?, password=?, wins=?, losses=?, draws=? WHERE id=?")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %v", err)
	}

	res, err := stmt.Exec(user.Name, user.Email, user.Password, user.Wins, user.Losses, user.Draws, id)
	if err != nil {
		return 0, fmt.Errorf("updateUserById: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("updateUserById: %v", err)
	}

	return rowsAffected, nil
}
