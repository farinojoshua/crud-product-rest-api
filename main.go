package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Products struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"proce"`
}

var database = map[int]Products{}

// id untuk enumearasi di db
var lastID = 0

func main() {
	// 1. buat route multiplexer
	mux := http.NewServeMux()

	// 3. tambahkan
	mux.HandleFunc("GET /products", listProduct)
	mux.HandleFunc("POST /products", createProduct)
	mux.HandleFunc("PUT /products/{id}", updateProduct)
	mux.HandleFunc("DELETE /products/{id}", deleteProduct)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}

// 2. fungsi handler
func listProduct(w http.ResponseWriter, r *http.Request) {
	// slice untuk response
	var products []Products

	// melakukan iterasi pada map database untuk memasukkan nilai ke structnya
	for _, v := range database {
		products = append(products, v)
	}

	// ubah menjadi json
	data, err := json.Marshal(products)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("Terjadi Kesalahan"))
	}

	// hasil dari json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}

	var products Products
	err = json.Unmarshal(bodyByte, &products)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}

	// inkrement nomor urut
	lastID++

	products.ID = lastID

	database[products.ID] = products

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte("Request berhasil ditambahkan"))
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}

	var products Products
	err = json.Unmarshal(bodyByte, &products)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Kesalahan dalam request"))
	}

	// supaya ID terbaru tidak diganti dengan id baru
	products.ID = productIDInt

	database[productIDInt] = products
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")

	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte("Kesalahan dalam request"))
	}

	delete(database, productIDInt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
	w.Write([]byte("Berhasil menghapus data"))
}
