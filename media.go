package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UploadResult struct {
	Media []struct {
		Link  string
		Title string
	}
}

func upload_media(filename string) {
	var ur UploadResult

	j := getApiFetcher("media/new")
	j.Files["media"] = filename
	resp, err := j.Method("POST").Send()

	if err != nil {
		log.Fatalln(">>Error: ", err)
	}

	if err := json.Unmarshal(resp.Bytes, &ur); err != nil {
		log.Fatal("Error parsing:", err)
	}

	fmt.Println(ur.Media[0].Link)
}
