package main

import (
	"fmt"
	"net/http"
)

// main is the entry point of our Go application.
// It will eventually run the web server.
func main() {
	fmt.Println("Hello from PeriFyGo!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Simple test handler
		w.Write([]byte("PeriFyGo says hi!"))
	})

	// Let's run on port 8080
	http.ListenAndServe(":8080", nil)
}
