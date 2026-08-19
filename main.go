package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"time"

	// "log"
	"net/http"

	"github.com/MaminirinaEdwino/turbojwt"
)



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
		ExpDate:       expDate,
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

func UnixToFormattedDate(unixTimestamp int64) string {
	// 1. Convertir le timestamp Unix en type time.Time
	t := time.Unix(unixTimestamp, 0)

	// 2. Formater la date avec le pattern de référence de Go (2006-01-02 15:04:05)
	return t.Format("02/01/2006 à 15:04:05")
}

func handleVerifyTokenForm(w http.ResponseWriter, r *http.Request) {
	token := r.PostFormValue("token")
	value, err := turbojwt.Verify("secret", token, 0)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	tmpl, err := template.ParseFiles("views/verify_result.html")
	if err != nil {
		log.Fatalf("Erreur lors du parsing du template : %v", err)
	}
	float64ext, ok := value["exp"].(float64)

	if ok {
		fmt.Println(value["exp"])
	}
	int64iat, _ := value["iat"].(float64)
	// Exemple 1: Avec une Map (directement issue du JSON)
	dataMap := map[string]interface{}{
		"email":         value["email"],
		"exp":           UnixToFormattedDate(int64(float64ext)),
		"full_name":     value["full_name"],
		"iat":           UnixToFormattedDate(int64(int64iat)),
		"number_of_day": value["number_of_day"],
		"subscription":  value["subscription"],
	}

	// Exécution du template (écrit dans la sortie standard os.Stdout pour tester)
	log.Println("--- Rendu du template avec Map ---")
	err = tmpl.Execute(w, dataMap)
	if err != nil {
		log.Fatalf("Erreur lors de l'exécution du template : %v", err)
	}
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
	mux.HandleFunc("GET /", handleGeneratingTokenInterface)
	mux.HandleFunc("GET /verify", handleVerifyTolenInterface)
	mux.HandleFunc("POST /verify", handleVerifyToken)
	mux.HandleFunc("POST /verify-form", handleVerifyTokenForm)
	mux.HandleFunc("POST /token", handleGeneratingToken)
	startServer(mux, PORT)
}
