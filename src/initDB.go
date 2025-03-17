package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func initDB(db *sql.DB) {
	db.Exec(`
		CREATE TABLE IF NOT EXISTS dest (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			code VARCHAR(255) UNIQUE
		);
	`)
	db.Exec(`
		CREATE TABLE IF NOT EXISTS route (
			id SERIAL PRIMARY KEY,
			from_code VARCHAR(255),
			to_code VARCHAR(255),
			duration INTEGER,
			route VARCHAR(255),
			transport_type VARCHAR(255)
		);
		`)
	fmt.Println("Databases created!")

	var count int
	db.QueryRow("SELECT COUNT(*) FROM dest").Scan(&count)

	if count == 0 {
		addDests(db)
	} else {
		fmt.Printf("Already have %d destinations\n", count)
	}

	db.QueryRow("SELECT COUNT(*) FROM route").Scan(&count)

	if count == 0 {
		addRoutes(db)
	} else {
		fmt.Printf("Already have %d routes\n", count)
	}

	cache = &Cache{codeName: make(map[string]string), nameCode: make(map[string]string)}
	rows, _ := db.Query("SELECT code, name FROM dest")
	for rows.Next() {
		var code string
		var name string
		rows.Scan(&code, &name)
		cache.codeName[code] = name
		cache.nameCode[name] = code
	}

}

func addDests(db *sql.DB) {
	for _, country := range allCountries.Countries {
		if country.Title == "Россия" {
			for _, region := range country.Regions {
				for _, settlement := range region.Settlements {
					for _, station := range settlement.Stations {
						if station.TransportType == "plane" && settlement.Title != "" {
							db.Exec("INSERT INTO dest (name, code) VALUES ($1, $2)", settlement.Title, settlement.Codes.YandexCode)

						}
					}
				}
			}
		}
	}
	var count int
	db.QueryRow("SELECT COUNT(*) FROM dest").Scan(&count)
	fmt.Printf("Added %d destinations\n", count)
}

func addRoutes(db *sql.DB) {

	city_codes := make([]string, 0)
	rows, _ := db.Query("SELECT code FROM dest")
	for rows.Next() {
		var code string
		rows.Scan(&code)
		city_codes = append(city_codes, code)
	}

	for i := 0; i < len(city_codes); i++ {
		for j := i + 1; j < len(city_codes); j++ {
			from := city_codes[i]
			to := city_codes[j]

			req, _ := http.NewRequest("GET", "https://api.rasp.yandex.net/v3.0/search/?from="+from+"&to="+to, nil)

			req.Header.Add("Authorization", apiKey)

			client := &http.Client{}
			resp, _ := client.Do(req)

			body, _ := ioutil.ReadAll(resp.Body)

			var result Routes
			json.Unmarshal(body, &result)

			for _, route := range result.Route {
				db.Exec("INSERT INTO route (from_code, to_code, duration, route, transport_type) VALUES ($1, $2, $3, $4, $5)", from, to, route.Duration, route.Thread.Number, route.From.TransportType)
			}
		}
		var count int
		db.QueryRow("SELECT COUNT(*) FROM route").Scan(&count)
		fmt.Printf("Scanned %d cities, found %d routes\n", i+1, count)
	}
}
