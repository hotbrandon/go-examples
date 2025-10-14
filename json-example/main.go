package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Todo struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func example1() {
	url := "https://jsonplaceholder.typicode.com/todos/2"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	// The http.Get() function will only return an error for network-level
	// issues (like DNS resolution failures, connection timeouts, etc.),
	// but it won't return an error for HTTP error status codes like 404, 500, etc.
	// Those are considered "successful" HTTP responses from a networking perspective.
	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("before unmarshalling: \n%s", string(body))
		}
		var todo Todo
		err = json.Unmarshal(body, &todo)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("after unmarshalling: %v", todo)
	} else {
		log.Println("response status code:", response.StatusCode)
	}
}

func example2() {
	url := "https://jsonplaceholder.typicode.com/todos/2"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	// The http.Get() function will only return an error for network-level
	// issues (like DNS resolution failures, connection timeouts, etc.),
	// but it won't return an error for HTTP error status codes like 404, 500, etc.
	// Those are considered "successful" HTTP responses from a networking perspective.
	if response.StatusCode == http.StatusOK {

		var todo Todo
		// fewer lines of code with json.NewDecoder
		dec := json.NewDecoder(response.Body)
		if err := dec.Decode(&todo); err != nil {
			// Decode returns an error if JSON is invalid or incomplete
			log.Println("error decoding JSON:", err)
			return
		}

		log.Printf("after unmarshalling: %v", todo)
	} else {
		log.Println("response status code:", response.StatusCode)
	}
}

func main() {
	// tutorial: https://www.youtube.com/watch?v=Vr63uGL7NrU
	example2()

}
