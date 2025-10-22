package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server on :9000")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello world!")
	})

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		return
	}
}
