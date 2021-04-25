package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"github.com/pkg/errors"
)

var db *sql.DB

type User struct {
	ID   int64
	Name string
	Age  int
}

func GetUserByUserID(ctx context.Context, userID int64) (*User, error) {
	u := &User{}
	query := "SELECT id, name, age FROM user WHERE id=?"
	err := db.QueryRowContext(ctx, query, userID).Scan(&u.ID, &u.Name, &u.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithMessage(err, query)
		}
		return nil, errors.Wrap(err, "get users by user_id")
	}
	return u, nil
}

func init() {
	os.Remove("users.db")
	log.Println("create user")
	file, err := os.Create("users.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalln(err.Error())
	}
	//defer db.Close()
	createTable(db)
	insertUser(db, "Yang", 20, "F")
}

func createTable(db *sql.DB) {
	createUserSQL := `CREATE TABLE user (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT,
		"age" integer,
		"sex" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create user table...")
	statement, err := db.Prepare(createUserSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("user table created")
}

func insertUser(db *sql.DB, name string, age int, sex string) {
	insertSQL := `INSERT INTO user (name, age, sex) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(name, age, sex)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
