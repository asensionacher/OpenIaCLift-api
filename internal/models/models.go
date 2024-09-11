package models

type User struct {
	Username string `json:"username"`
	Hash     string `json:"hash"`
}

type Login struct {
	Username     string `json:"username"`
	FlatPassword string `json:"flatpassword"`
}
