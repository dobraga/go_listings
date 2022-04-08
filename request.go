package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var client http.Client = http.Client{}

func MakeRequest(
	url string,
	query map[string]interface{},
	headers map[string]string,
) []byte {
	req, err := http.NewRequest("GET", url, nil)
	Check(err)

	// Query String
	q := req.URL.Query()
	for key, element := range query {
		q.Add(key, fmt.Sprintf("%v", element))
	}

	req.URL.RawQuery = q.Encode()

	// Headers
	for key, element := range headers {
		req.Header.Add(key, fmt.Sprintf("%v", element))
	}

	// Request
	resp, err := client.Do(req)
	Check(err)
	defer resp.Body.Close()

	// Response to interface
	bytes_data, err := ioutil.ReadAll(resp.Body)
	Check(err)

	return bytes_data
}
