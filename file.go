package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
	"log"
)

type InputRequest struct {
	ToSort [][]int `json:"to_sort"`
}

type OutputResponse struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs int64    `json:"time_ns"`
}

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)

	port := getPort()

	// Starting the server
	log.Printf("Server listening on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		// Default to 8000 if PORT environment variable is not set
		port = "8000"
	}
	return port
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, false)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, true)
}

func handleRequest(w http.ResponseWriter, r *http.Request, concurrent bool) {
	var req InputRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	var sortedArrays [][]int

	if concurrent {
		sortedArrays = sortConcurrent(req.ToSort)
	} else {
		sortedArrays = sortSequential(req.ToSort)
	}
	timeTaken := time.Since(startTime)

	response := OutputResponse{
		SortedArrays: sortedArrays,
		TimeNs:       timeTaken.Nanoseconds(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sortSequential(arrays [][]int) [][]int {
	sortedArrays := make([][]int, len(arrays))

	for i, arr := range arrays {
		sortedArr := make([]int, len(arr))
		copy(sortedArr, arr)
		sort.Ints(sortedArr)
		sortedArrays[i] = sortedArr
	}
	return sortedArrays
}

func sortConcurrent(arrays [][]int) [][]int {
	var wg sync.WaitGroup
	wg.Add(len(arrays))

	sortedArrays := make([][]int, len(arrays))

	for i, arr := range arrays {
		go func(i int, arr []int) {
			defer wg.Done()
			sortedArr := make([]int, len(arr))
			copy(sortedArr, arr)
			sort.Ints(sortedArr)
			sortedArrays[i] = sortedArr
		}(i, arr)
	}
	wg.Wait()

	return sortedArrays
}
