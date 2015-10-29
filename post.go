package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// struct for file be reading, create a new post
// slightly different with category and tags
type Page struct {
	Title, Content, Category, Status, Tags string
	Date                                   time.Time
}

// create new post
func do_post(filename string) {

	page := readParseFile(filename)

	j := getApiFetcher("posts/new")
	j.Params.Add("title", page.Title)
	j.Params.Add("date", page.Date.Format(time.RFC3339))
	j.Params.Add("content", page.Content)
	j.Params.Add("status", page.Status)
	j.Params.Add("categories", page.Category)
	j.Params.Add("publicize", "0")
	j.Params.Add("tags", page.Tags)

	resp, err := j.Method("POST").Send()
	if err != nil {
		log.Fatalln(">>Error: ", err)
	}

	newurl := parseNewPostResponse(resp.Bytes)
	fmt.Println("New Post:", newurl)
}

// read and parse markdown filename
func readParseFile(filename string) (page Page) {

	// setup default page struct
	page = Page{
		Title:    "",
		Content:  "",
		Category: "",
		Date:     time.Now(),
		Tags:     "",
		Status:   "publish",
	}

	var data, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(">>Error: can't read file:", filename)
	}

	// parse front matter from --- to ---
	var lines = strings.Split(string(data), "\n")
	var found = 0
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if found == 1 {
			// parse line for param
			colonIndex := strings.Index(line, ":")
			if colonIndex > 0 {
				key := strings.TrimSpace(line[:colonIndex])
				value := strings.TrimSpace(line[colonIndex+1:])
				value = strings.Trim(value, "\"") //remove quotes
				switch key {
				case "title":
					page.Title = value
				case "date":
					page.Date, _ = time.Parse("2006-01-02", value)
				case "category":
					page.Category = value
				case "tags":
					page.Tags = value
				case "status":
					page.Status = value
				}
			}
		} else if found >= 2 {
			// params over
			lines = lines[i:]
			break
		}

		if line == "---" {
			found += 1
		}
	}

	// slurp rest of content
	page.Content = strings.Join(lines, "\n")
	return page
}

// extract URL from json response data of new post
func parseNewPostResponse(data []byte) string {

	var rs struct{ Url string }

	if err := json.Unmarshal(data, &rs); err != nil {
		log.Fatalf("Error parsing: {} \n\n {}", data, err)
	}

	return rs.Url
}
