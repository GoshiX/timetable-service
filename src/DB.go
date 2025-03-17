package main

import "database/sql"

type Edge struct {
	From      string
	To        string
	Duration  int
	Route     string
	Transport string
}

type Graph struct {
	Edges map[string][]Edge
}

type Cache struct {
	codeName map[string]string
	nameCode map[string]string
}

func getCityCode(db *sql.DB, city string) string {
	return cache.nameCode[city]
}

func getCityName(db *sql.DB, code string) string {
	return cache.codeName[code]
}

func getGraph(db *sql.DB) *Graph {
	graph := &Graph{}
	graph.Edges = make(map[string][]Edge)
	rows, _ := db.Query("SELECT from_code, to_code, duration, route, transport_type FROM route")
	for rows.Next() {
		var from string
		var to string
		var duration int
		var route string
		var transport string
		rows.Scan(&from, &to, &duration, &route, &transport)
		if _, ok := graph.Edges[from]; !ok {
			graph.Edges[from] = make([]Edge, 0)
		}
		graph.Edges[from] = append(graph.Edges[from], Edge{From: from, To: to, Duration: duration, Route: route, Transport: transport})
		if _, ok := graph.Edges[to]; !ok {
			graph.Edges[to] = make([]Edge, 0)
		}
		graph.Edges[to] = append(graph.Edges[to], Edge{From: to, To: from, Duration: duration, Route: route, Transport: transport})
	}
	return graph
}
