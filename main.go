package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MaminirinaEdwino/turbojwt"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "home",
	})
}

type bodyType struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	NumberOfDay  int    `json:"number_of_day"`
	PaymentRef   string `json:"payment_ref"`
	Subscription string `json:"subscription"`
}

func handleGeneratingToken(w http.ResponseWriter, r *http.Request) {
	var bodyValue bodyType
	json.NewDecoder(r.Body).Decode(&bodyValue)
	token, err := turbojwt.Encode("secret", map[string]any{
		"full_name":     bodyValue.FullName,
		"email":         bodyValue.Email,
		"number_of_day": bodyValue.NumberOfDay,
		"payment_ref":   bodyValue.PaymentRef,
		"subscription":  bodyValue.Subscription,
	}, float64(bodyValue.NumberOfDay*24), 0)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(token)
}

func handleGeneratingTokenInterface(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
func handleVerifyToken(w http.ResponseWriter, r *http.Request) {
	var token map[string]string
	json.NewDecoder(r.Body).Decode(&token)
	value, err := turbojwt.Verify("secret", token["token"], 0)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(value)
}
func handleVerifyTolenInterface(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "verify.html")
}
func startServer(mux *http.ServeMux, PORT int) {
	fmt.Printf("Server started on port : http://localhost:%d\n", PORT)
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
	mux.HandleFunc("GET /token", handleGeneratingTokenInterface)
	mux.HandleFunc("POST /token", handleGeneratingToken)
	startServer(mux, PORT)
}
