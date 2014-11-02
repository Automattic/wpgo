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

type StreakResult struct {
	Count	int			`json:"found"`
	Posts	[]Page		`json:"posts"`
}

type TopPostsResult struct {
	Day      string `json:"date"`
	TopPosts []struct {
		PostId int
		Title  string
		Views  int
	} `json:"top-posts"`
}

func get_stats(stat_type string) {
	if stat_type == "streak" {
		getStreakData()
	} else {
		getSummary()
	}
}

// This returns a stats summary view of view, visitors for today and yesterday.
// Also shows top posts and number of views.
func getSummary() {
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

// This returns JSON data with a day timestamp and number of posts for that day.
// This format works with cal-heatmap library
func getStreakData() {
	// fetch posts after Jan 2014 (limit 100)
	// after: DATE
	// number: 100
	// fields: ID,title,date
	f, url := get_api_fetcher("posts")
	f.Params.Add("after", "2014-01-01")
	f.Params.Add("number", "100")
	f.Params.Add("fields", "date")
	result, err := f.Fetch(url, "GET")
	if err != nil {
		log.Fatalln(">>Error Fetching: ", err)
	}

	var s StreakResult
	if err := json.Unmarshal([]byte(result), &s); err != nil {
		log.Fatalln("Error parsing:", err)
	}

	// TODO if found == 100, need to query again

	// TODO if multiple posts on single day, count
	// encode as real JSON
	jsmap := make(map[string]int)
	for _, p := range s.Posts {
		ds := fmt.Sprintf("%v", p.Date.Unix())
		jsmap[ds] = 1
	}
	b, err := json.Marshal(jsmap)
	if err != nil {
		log.Fatalln("Error Marshaling:", err)
	}
	fmt.Print(string(b))

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
