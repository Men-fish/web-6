package main

import (
	"fmt"
	"log"
	"net/http"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Parameter 'name' is required", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("Hello, %s!", name)
	w.Write([]byte(response))
}

func main() {
	http.HandleFunc("/api/user", userHandler)

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// curl http://localhost:8080/api/user?name=Men-fish
