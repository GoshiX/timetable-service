package main

import (
	"sort"
)

type resRoute struct {
	Part []Edge
}

func (r *resRoute) Duration() int {
	res := 0
	for _, e := range r.Part {
		res += e.Duration
	}
	return res
}

func findPath(graph *Graph, from string, to string) *[]resRoute {
	if from == to {
		return &[]resRoute{{Part: []Edge{}}}
	}
	res := make([]resRoute, 0)
	q := make([]resRoute, 0)
	q = append(q, resRoute{Part: []Edge{Edge{
		From: from,
		To:   from,
	}}})
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		for _, e := range graph.Edges[getCityCode(db, cur.Part[len(cur.Part)-1].To)] {
			newCur := resRoute{Part: make([]Edge, 0)}
			newCur.Part = append(newCur.Part, cur.Part...)
			if len(newCur.Part) > 1 {
				newCur.Part = append(newCur.Part, Edge{
					From:     "Transfer",
					To:       "Transfer",
					Duration: 45 * 60, // 45 min
				})
			}
			newCur.Part = append(newCur.Part, Edge{
				From:      getCityName(db, e.From),
				To:        getCityName(db, e.To),
				Duration:  e.Duration,
				Route:     e.Route,
				Transport: e.Transport,
			})
			if getCityName(db, e.To) == to {
				res = append(res, newCur)
				continue
			}
			if len(newCur.Part) < 4 {
				q = append(q, newCur)
			}
		}
	}
	for i := 0; i < len(res); i++ {
		res[i].Part = res[i].Part[1:]
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Duration() < res[j].Duration()
	})
	return &res
}
