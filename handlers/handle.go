package handlers

import (
	"fmt"
	"net/http"

	"database/sql"
)

var db *sql.DB

func SetDB(database *sql.DB) {
    db = database
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Home Page!")
}