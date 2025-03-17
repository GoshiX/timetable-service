package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Codes struct {
	YandexCode string `json:"yandex_code"`
	EsrCode    string `json:"esr_code,omitempty"` // omitempty означает, что поле может отсутствовать в JSON
}

// Station представляет собой структуру для хранения информации о станции.
type Station struct {
	Direction     string  `json:"direction"`
	Codes         Codes   `json:"codes"`
	StationType   string  `json:"station_type"`
	Title         string  `json:"title"`
	Longitude     float64 `json:"longitude"`
	TransportType string  `json:"transport_type"`
	Latitude      float64 `json:"latitude"`
}

// Settlement представляет собой структуру для хранения информации о населенном пункте.
type Settlement struct {
	Title      string    `json:"title"`
	Codes      Codes     `json:"codes"`
	Stations   []Station `json:"stations"`
	StationMap map[string]*Station
}

// Region представляет собой структуру для хранения информации о регионе.
type Region struct {
	Title          string       `json:"title"`
	Codes          Codes        `json:"codes"`
	Settlements    []Settlement `json:"settlements"`
	SettlementsMap map[string]*Settlement
}

// Country представляет собой структуру для хранения информации о стране.
type Country struct {
	Title      string   `json:"title"`
	Codes      Codes    `json:"codes"`
	Regions    []Region `json:"regions"`
	RegionsMap map[string]*Region
}

// Countries представляет собой структуру для хранения списка стран.
type AllCountries struct {
	Countries    []Country `json:"countries"`
	CountriesMap map[string]*Country
}

var (
	url string = "https://api.rasp.yandex.net/v3.0/stations_list/"
)

func getAllStationsList(apiKey string) string {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", apiKey)

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

func getPrettyData() *AllCountries {
	dataRaw, _ := os.ReadFile("data/allStations.json")
	var data *AllCountries = &AllCountries{}
	json.Unmarshal(dataRaw, data)
	return data
}

type Routes struct {
	Route []Route `json:"segments"`
}

type Route struct {
	From     Station `json:"from"`
	To       Station `json:"to"`
	Thread   Thread  `json:"thread"`
	Duration float32 `json:"duration"`
}

type Thread struct {
	Number  string `json:"number"`
	Vehicle string `json:"vehicle"`
}

func parseRoutes() {
	dataRaw, _ := os.ReadFile("tmp.json")
	var result Routes
	json.Unmarshal(dataRaw, &result)

	for _, route := range result.Route {
		fmt.Println(route.From.Title, route.To.Title, route.Duration, route.Thread.Vehicle, route.Thread.Number)
	}
}
