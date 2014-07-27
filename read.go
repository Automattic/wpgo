package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type author struct {
	ID          int
	Login, Name string
}

type Post struct {
	Id      int       `json:"ID"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	URL     string    `json:"URL"`
	Author  author    `json:"author"`
	Content string    `json:"content"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}


func get_latest() {
	posts := parseFetchPosts()
	for _, post := range posts {
		fmt.Printf("[%d] %-55s (%s)\n", post.Id, post.Title, post.Author.Name)
	}
}

// fetch single post
func get_single_post(post_id string) {
	post := parseFetchPost(post_id)
	fmt.Println(post.Title)
	fmt.Println(scrub_html(post.Content))
	fmt.Println(post.URL)
}

// fetch and parse list of posts
func parseFetchPosts() ([]Post) {
	f, url := get_api_fetcher("posts/")
	result, err := f.Fetch(url, "GET")
	if err != nil {
		log.Fatalln(">>Error: ", err)
	}

	var h Posts
	if err := json.Unmarshal([]byte(result), &h); err != nil {
		log.Fatal("Error parsing:", err)
	}

	return h.Posts
}

// parse single post
func parseFetchPost(post_id string) (p Post) {
	f, url := get_api_fetcher("posts/" + post_id)
	result, err := f.Fetch(url, "GET")
	if err != nil {
		log.Fatalln(">>Error: ", err)
	}

	if err := json.Unmarshal([]byte(result), &p); err != nil {
		log.Fatal("Error parsing:", err)
	}
	return p
}
