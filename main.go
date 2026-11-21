package main 

import (
	"log"
	"net/http"

	"together/database"
	"together/handlers"

	
)
 
func main() {

	db := database.ConnectDB()
    defer db.Close()

    handlers.SetDB(db)


	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.Home)

	mux.HandleFunc("GET /space", handlers.GetSpace)

	mux.HandleFunc("POST /create-space", handlers.CreateSpace)

	// mux.HandleFunc("POST /enter-space", )

	// mux.HandleFunc("POST /add-media", )

	// mux.HandleFunc("POST /sync-media", )

	log.Println("Server running on :8081")

	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}

}	