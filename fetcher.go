//
// A library to make it a bit easier to do HTTP fetches
// supports adding headers, posting forms, parameters
// and uploading files
//

package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Fetcher struct {
	Params, Header, Files map[string]string
}

// NewFetcher creates a fetcher request instance
func NewFetcher() (f Fetcher) {
	f.Params = map[string]string{}
	f.Header = map[string]string{}
	f.Files = map[string]string{}
	return f
}

// Fetch executes the fetcher request
func (f Fetcher) Fetch(url, method string) (result string, err error) {

	var reqBody io.Reader
	var contentType string

	// check if post and add post params
	if method == "POST" {
		reqBody, contentType, err = f.createPostBody()
		if err != nil {
			return "", err
		}
	} else {
		method = "GET"
	}

	// build request object
	client := &http.Client{}
	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return "", err
	}

	// need to add header to request for content-type
	// this sets boundaries and builds proper header type
	if method == "POST" {
		request.Header.Add("Content-Type", contentType)
	}

	// add additional user header values
	for k, v := range f.Header {
		request.Header.Add(k, v)
	}

	// execute request
	res, err := client.Do(request)
	if err != nil {
		return "", err
	}

	// process response
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	result = string(body)
	return
}

// create body for post - includes files, params
func (f Fetcher) createPostBody() (body io.Reader, contentType string, err error) {

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// add files if we are uploading a file
	for k, v := range f.Files {
		file, err := os.Open(v)
		if err != nil {
			return nil, "", err
		}

		part, err := writer.CreateFormFile(k, filepath.Base(v))
		if err != nil {
			return nil, "", err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, "", err
		}
	}

	// add parameters if there are parameters
	for k, v := range f.Params {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close()
	if err != nil {
		return
	}

	// content type might be different due to file uploads
	contentType = writer.FormDataContentType()
	body = &b
	return body, contentType, nil

}
