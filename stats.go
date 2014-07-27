package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type StatResult struct {
	Day   string `json:"day"`
	Stats struct {
		VisitorsToday     int `json:"visitors_today"`
		VisitorsYesterday int `json:"visitors_yesterday"`
		ViewsToday        int `json:"views_today"`
		ViewsYesterday    int `json:"views_yesterday"`
	}
}

type TopPostsResult struct {
	Day      string `json:"date"`
	TopPosts []struct {
		PostId int
		Title  string
		Views  int
	} `json:"top-posts"`
}

func get_stats() {
	s := parseFetchStats()
	fmt.Println("---------------------------------------------")
	fmt.Println("       TODAY          |        YESTERDAY")
	fmt.Printf("  %4d     %4d", s.Stats.VisitorsToday, s.Stats.ViewsToday)
	fmt.Printf("       |      %4d     %4d \n", s.Stats.VisitorsYesterday, s.Stats.ViewsYesterday)
	fmt.Println(" Visitors  Views      |     Visitors  Views ")
	fmt.Println("---------------------------------------------")
	fmt.Println(" TOP POSTS")
	fmt.Println("---------------------------------------------")
	top_posts := parseFetchTopPosts(10)
	for _, tp := range top_posts.TopPosts {
		if len(tp.Title) > 36 {
			tp.Title = tp.Title[:33] + "..."
		}
		fmt.Printf("%-36s %6d \n", tp.Title, tp.Views)
	}
	fmt.Println("---------------------------------------------")
}

func parseFetchStats() (s StatResult) {

	f, url := get_api_fetcher("stats")
	result, err := f.Fetch(url, "GET")
	if err != nil {
		log.Fatalln(">>Error: ", err)
	}

	if err := json.Unmarshal([]byte(result), &s); err != nil {
		log.Fatal("Error parsing:", err)
	}

	return s
}

func parseFetchTopPosts(limit int) (tp TopPostsResult) {
	f, url := get_api_fetcher("stats/top-posts")
	result, err := f.Fetch(url, "GET")
	if err != nil {
		log.Fatalln(">>Error: ", err)
	}

	if err := json.Unmarshal([]byte(result), &tp); err != nil {
		log.Fatal("Error parsing:", err)
	}

	if len(tp.TopPosts) > limit {
		tp.TopPosts = tp.TopPosts[:limit]
	}
	return tp

}
