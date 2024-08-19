package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    log "github.com/go-kit/kit/log"
    "math/rand"
    "net/http"
    "os"
    "strconv"
    "time"
)

var logger log.Logger

type GithubOrg struct {
	OrgName   	string    
    Token     	string    
	CreatedAt	time.Time 
}

type Repo struct {
	Name      	string    
    GithubOrg	GithubOrg 
    CreatedAt 	time.Time 
}

type PullRequest struct {
	Id			integer    
    Repo		Repo
    CreatedAt 	time.Time 
	// Reviewed
	// CI executed
	// etc
}

// https://github.com/adityajoshi12/restapi-golang/blob/master/main.go
func main() {
    r := mux.NewRouter()

    r.HandleFunc("/repo", NewRepo).Methods("POST")
    http.Handle("/", r)
    if err := http.ListenAndServe(":3000", r); err != nil {
        logger.Log("status", "fatal", "err", err)
        os.Exit(1)
    }
}

func NewOrg(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

func NewRepo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    json.NewEncoder(w).Encode(posts)
}