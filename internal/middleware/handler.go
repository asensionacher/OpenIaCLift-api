package middleware

import (
	"OpenIaCLiftAPI/internal/models"
	"encoding/json"
	"fmt"
	"net/http"

	"database/sql" // package to encode and decode the json into struct and vice versa

	"log" // used to access the request and response object of the api

	"os" // used to read the environment variable

	// package used to covert string into int type
	// used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// LoginUser checks if password is correct. If user not exists creates user with default password
func LoginUser(w http.ResponseWriter, r *http.Request) {
	// create an empty user of type models.User
	var login models.Login

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// Check if user exists. If not, create one admin - admin and return the object.
	user, _ := checkUserExists()

	// Check if login is correct
	loginCorrect := checkPassword(user.Hash, login.FlatPassword)

	// format a response object
	res := response{
		Message: "OK",
	}
	if !loginCorrect {
		res = response{
			Message: "NOOK",
		}
	}
	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------

func checkUserExists() (models.User, error) {
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a user of models.User type
	var user models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users WHERE username=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, "admin")
	// unmarshal the row object to user
	err := row.Scan(&user.Username, &user.Hash)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned, create user")
		sqlStatement := `INSERT INTO users (username, hash) VALUES ($1, $2) RETURNING username`
		hashedPassword, _ := hashPassword("admin")
		var username string
		db.QueryRow(sqlStatement, "admin", hashedPassword).Scan(&username)
		user.Username = username
		user.Hash = hashedPassword
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return user, err
}
