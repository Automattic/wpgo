package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// upload a single file
func upload_media(filename string) {
	var ur struct {
		Media []struct {
			Link  string
			Title string
		}
	}

	j := getApiFetcher("media/new")
	j.Files["media[]"] = filename
	resp, err := j.Method("POST").Send()

	if err != nil {
		log.Fatalln(">>Error: ", err)
	}

	if err := json.Unmarshal(resp.Bytes, &ur); err != nil {
		log.Fatal("Error parsing:", err)
	}

	if len(ur.Media) > 0 {
		fmt.Println(ur.Media[0].Link)
	} else {
		fmt.Println("Error: No link in results")
		fmt.Println(resp.StatusCode)
	}

}
