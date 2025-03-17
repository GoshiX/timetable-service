package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var (
	apiKey       string
	allCountries = getPrettyData()
	db           *sql.DB
	g            *Graph
	cache        *Cache
)

func main() {

	apiKey = os.Getenv("API_KEY")
	psqlInfo := os.Getenv("DATABASE_URL")
	time.Sleep(time.Second)
	db, _ = sql.Open("postgres", psqlInfo)
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	fmt.Println("Successfully connected to the database!")

	initDB(db)
	g = getGraph(db)

	http.HandleFunc("/available_dest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		rows, _ := db.Query("SELECT name FROM dest")
		defer rows.Close()

		destinations := make([]string, 0)
		for rows.Next() {
			var name string
			rows.Scan(&name)
			destinations = append(destinations, name)
		}

		json.NewEncoder(w).Encode(destinations)
	})

	http.HandleFunc("/route", routeHandler)

	fmt.Printf("Starting server on :8000...\n")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad city name"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res := findPath(g, from, to)
	json.NewEncoder(w).Encode(res)

}
