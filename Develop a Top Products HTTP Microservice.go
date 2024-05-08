package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	testServerURL = "http://20.244.56.144/test/companies/"
)

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Rating   float64 `json:"rating"`
	Discount float64 `json:"discount"`
	Company  string  `json:"company"`
}

type ProductResponse struct {
	Products []Product `json:"products"`
}

func main() {
	http.HandleFunc("/categories/", productsHandler)
	http.HandleFunc("/categories/", productDetailsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := r.URL.Path[len("/categories/"):]
	// Extract query parameters
	nStr := r.URL.Query().Get("n")
	pageStr := r.URL.Query().Get("page")
	sortBy := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")

	// Parse query parameters
	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, "Invalid value for 'n'", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1 // Default page number
	}

	// Fetch products from test server
	products := fetchProducts(categoryName, n, page, sortBy, sortOrder)

	// Encode products as JSON and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func productDetailsHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := r.URL.Path[len("/categories/"):]
	productID := r.URL.Path[len("/categories/"+categoryName+"/products/"):]

	// Fetch product details from test server
	product := fetchProductDetails(categoryName, productID)

	// Encode product details as JSON and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func fetchProducts(categoryName string, n, page int, sortBy, sortOrder string) ProductResponse {
	// Construct URL for fetching products from test server
	url := fmt.Sprintf("%s/categories/%s/products?top=%d&page=%d&sort=%s&order=%s", testServerURL, categoryName, n, page, sortBy, sortOrder)

	// Make GET request to test server
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching products: %v\n", err)
		return ProductResponse{}
	}
	defer resp.Body.Close()

	// Decode response
	var productResponse ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&productResponse); err != nil {
		log.Printf("Error decoding product response: %v\n", err)
		return ProductResponse{}
	}

	return productResponse
}

func fetchProductDetails(categoryName, productID string) Product {
	// Construct URL for fetching product details from test server
	url := fmt.Sprintf("%s/categories/%s/products/%s", testServerURL, categoryName, productID)

	// Make GET request to test server
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching product details: %v\n", err)
		return Product{}
	}
	defer resp.Body.Close()

	// Decode response
	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		log.Printf("Error decoding product details response: %v\n", err)
		return Product{}
	}

	return product
}
