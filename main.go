package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "home",
	})
}

type bodyType struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	NumberOfDay int    `json:"number_of_day"`
	PaymentRef  string `json:"payment_ref"`
}

func handleGeneratingToken(w http.ResponseWriter, r *http.Request) {
	var bodyValue bodyType
	json.NewDecoder(r.Body).Decode(&bodyValue)
	json.NewEncoder(w).Encode(bodyValue)
}

func startServer(mux *http.ServeMux, PORT int) {
	fmt.Printf("Server started on port : %d\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
	if err != nil {
		fmt.Println(err)
		startServer(mux, PORT+1)
	}
}
func main() {
	mux := http.NewServeMux()
	PORT := 1627
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("POST /token", handleGeneratingToken)
	startServer(mux, PORT)
}
