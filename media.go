package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UploadResult struct {
	Media []struct {
		Link string
		Title string
	}
}


func upload_media(filename string) {
	var ur UploadResult
	
	f, url := get_api_fetcher("media/new")
	f.Files["media"] = filename
	result, err := f.Fetch(url, "POST")

	if err != nil {
		log.Fatalln(">>Error: ", err)
	}

	if err := json.Unmarshal([]byte(result), &ur); err != nil {
		log.Fatal("Error parsing:", err)
	}

	fmt.Println(ur.Media[0].Link)
}
