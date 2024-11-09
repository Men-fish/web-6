package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	counter int
	mu      sync.Mutex
)

func getCountHandler(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(w, "Current count: %d", counter)
}

func postCountHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	countValue, ok := data["count"]
	if !ok {
		http.Error(w, "Key 'count' not found", http.StatusBadRequest)
		return
	}

	countFloat, ok := countValue.(float64)
	if !ok {
		http.Error(w, "это не число", http.StatusBadRequest)
		return
	}
	countInt := int(countFloat)
	mu.Lock()
	counter += countInt
	mu.Unlock()

	fmt.Fprintf(w, "Count incremented by %d, new count: %d", countInt, counter)
}
func main() {
	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCountHandler(w, r)
		case http.MethodPost:
			postCountHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on http://localhost:3333")
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// curl -X POST -H "Content-Type: application/json" -d '{"count": 5}' http://localhost:3333/count
// curl http://localhost:3333/count
