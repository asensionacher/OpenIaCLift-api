// TODO: Estructura de carpetas (db en internal, setup, etc)
package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var logger log.Logger

type GithubOrg struct {
	Name      string    `json:"Name"`
	Token     string    `json:"Token"`
	CreatedAt time.Time `json:"CreatedAt,omitempty"`
}

type Repo struct {
	Id        int       `json:"Id"`
	Name      string    `json:"Name"`
	CreatedAt time.Time `json:"CreatedAt,omitempty"`
}

type PullRequest struct {
	Id        int       `json:"Id"`
	Repo      Repo      `json:"Repo"`
	CreatedAt time.Time `json:"CreatedAt,omitempty"`
	// Reviewed
	// CI executed
	// etc
}

// https://github.com/adityajoshi12/restapi-golang/blob/master/main.go
func main() {
	// Init MySQL
	InitDb()

	r := mux.NewRouter()

	r.HandleFunc("/org", NewOrg).Methods("POST")
	r.HandleFunc("/repo", NewRepo).Methods("POST")
	http.Handle("/", r)
	if err := http.ListenAndServe(":3000", r); err != nil {
		logger.Log("status", "fatal", "err", err)
		os.Exit(1)
	}
}

const DBNAME string = "OpenIaCLift"

var (
	db  *sql.DB
	cfg = mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DBAddress"),
		DBName: DBNAME,
	}
	dbConnectionString = cfg.User + ":" + cfg.Passwd + "@" + cfg.Net + "(" + cfg.Addr + ")/"
)

func InitDb() {
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", dbConnectionString)
	if err != nil {
		logger := log.NewLogfmtLogger(os.Stdout)
		level.Error(logger).Log("msg", err)
		os.Exit(-1)
	}
	_, err = db.Exec("CREATE DATABASE " + cfg.DBName)
	if err != nil {
		logger := log.NewLogfmtLogger(os.Stdout)
		level.Error(logger).Log("msg", err)
	}
	_, err = db.Exec("USE " + cfg.DBName)
	if err != nil {
		logger := log.NewLogfmtLogger(os.Stdout)
		level.Error(logger).Log("msg", err)
	}
	// TODO: Create all tables
	// TODO: Acordarse de id + autoincrement
	// TODO: Antes hacer un poco de dibujo
	_, err = db.Exec("CREATE TABLE githuborg ( name varchar(255), token varchar(255), created_at datetime )")
	if err != nil {
		logger := log.NewLogfmtLogger(os.Stdout)
		level.Error(logger).Log("msg", err)
	}

	defer db.Close()
}

func DatabaseExec(query string) error {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DBAddress"),
		DBName: DBNAME,
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger := log.NewLogfmtLogger(os.Stdout)
		level.Error(logger).Log("msg", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + DBNAME)
	if err != nil {
		logger := log.NewLogfmtLogger(os.Stdout)
		level.Error(logger).Log("msg", err)
		return os.NewSyscallError(err.Error(), err)
	}
	return nil
}

func NewOrg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var g GithubOrg
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func NewRepo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	// json.NewEncoder(w).Encode(posts)
}
