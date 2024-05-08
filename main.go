package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	testServerURL  = "http://20.244.56.144/test/"
	windowSize     = 10
	requestTimeout = 500 * time.Millisecond
)

var (
	numbers     []int
	windowMutex sync.Mutex
)

type TestServerResponse struct {
	Numbers []int `json:"numbers"`
}

type Response struct {
	Numbers         []int   `json:"numbers"`
	WindowPrevState []int   `json:"windowPrevState"`
	WindowCurrState []int   `json:"windowCurrState"`
	Avg             float64 `json:"avg"`
}

func main() {
	http.HandleFunc("/numbers/", numbersHandler)
	log.Fatal(http.ListenAndServe(":9876", nil))
}

func numbersHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the qualifier from the request URL
	qualifier := r.URL.Path[len("/numbers/"):]
	// Fetch numbers from the test server
	testServerResponse := fetchNumbers(qualifier)
	// Update window state and calculate average
	updateWindowState(testServerResponse.Numbers)
	// Prepare response
	resp := Response{
		Numbers:         testServerResponse.Numbers,
		WindowPrevState: getPreviousWindow(),
		WindowCurrState: getCurrentWindow(),
		Avg:             calculateAverage(),
	}
	// Encode response as JSON and send it
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func fetchNumbers(qualifier string) TestServerResponse {
	url := testServerURL + qualifier
	client := http.Client{
		Timeout: requestTimeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error fetching numbers from test server: %v\n", err)
		return TestServerResponse{}
	}
	defer resp.Body.Close()
	var testServerResponse TestServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&testServerResponse); err != nil {
		log.Printf("Error decoding response from test server: %v\n", err)
		return TestServerResponse{}
	}
	return testServerResponse
}

func updateWindowState(newNumbers []int) {
	windowMutex.Lock()
	defer windowMutex.Unlock()
	// Append new numbers to the window
	numbers = append(numbers, newNumbers...)
	// Trim numbers if the window size exceeds
	if len(numbers) > windowSize {
		numbers = numbers[len(numbers)-windowSize:]
	}
}

func getPreviousWindow() []int {
	windowMutex.Lock()
	defer windowMutex.Unlock()
	// If the current window is smaller than the window size,
	// there's no previous window
	if len(numbers) < windowSize {
		return []int{}
	}
	return numbers[:windowSize]
}

func getCurrentWindow() []int {
	windowMutex.Lock()
	defer windowMutex.Unlock()
	return numbers
}

func calculateAverage() float64 {
	windowMutex.Lock()
	defer windowMutex.Unlock()
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return float64(sum) / float64(len(numbers))
}
