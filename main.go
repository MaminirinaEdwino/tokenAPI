package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	// "log"
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
	temp, _ := template.ParseFiles("views/token.html")

	nbrD, _ := strconv.Atoi(r.PostFormValue("number_of_day"))
	token, err := turbojwt.Encode("secret", map[string]any{
		"full_name":     r.PostFormValue("full_name"),
		"email":         r.PostFormValue("email"),
		"number_of_day": r.PostFormValue("number_of_day"),
		"subscription":  r.PostFormValue("subscription"),
	}, float64(nbrD*24), 0)
	decoded, _ := turbojwt.Verify("secret", token, 0)
	exp := decoded["exp"].(float64)
	expDate := time.Unix(int64(exp), 0).Format(time.UnixDate)
	data := struct {
		Full_name     string
		Email         string
		Number_of_day int
		Subscription  string
		Payment_ref   string
		Token         string
		ExpDate       string
	}{
		Full_name:     r.PostFormValue("full_name"),
		Email:         r.PostFormValue("email"),
		Number_of_day: nbrD,
		Token:         token,
		Subscription:  r.PostFormValue("subscription"),
		ExpDate: expDate,
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	temp.Execute(w, data)
}

func handleGeneratingTokenInterface(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

func handleVerifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.PostFormValue("token")
	value, err := turbojwt.Verify("secret", token, 0)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(value)
}

func handleVerifyTolenInterface(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/verify.html")
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

	styleServer := http.FileServer(http.Dir("styles"))

	PORT := 1627
	mux.Handle("GET /styles/", http.StripPrefix("/styles/", styleServer))
	// mux.HandleFunc("GET /", handleGeneratingTokenInterface)
	mux.HandleFunc("GET /token", handleGeneratingTokenInterface)
	mux.HandleFunc("GET /verify", handleVerifyTolenInterface)
	mux.HandleFunc("POST /verify", handleVerifyToken)
	mux.HandleFunc("POST /token", handleGeneratingToken)
	startServer(mux, PORT)
}
