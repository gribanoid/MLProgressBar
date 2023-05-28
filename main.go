package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var port int
var host string
var ml string

func init() {
	flag.IntVar(&port, "p", 8082, "application internal port")
	flag.StringVar(&host, "h", "http://localhost", "application host")
	flag.StringVar(&ml, "ml", "http://localhost:8090", "ml server")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/match/", match).Methods("POST")
	r.HandleFunc("/score/{score}/", score).Methods("GET")
	log.Println("server started")
	log.Fatalf("server stoped: %v", http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

func match(w http.ResponseWriter, r *http.Request) {
	var model Model
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, "parse body", http.StatusBadRequest)
		return
	}
	progressBar, err := CalculateMatch(ml, model)
	if err != nil {
		log.Printf("calculating match: %v", err)
		http.Error(w, "ml tool broken", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("%s:%d/score/%d/", host, port, progressBar), http.StatusFound)
}

func score(w http.ResponseWriter, r *http.Request) {
	progressBar, ok := mux.Vars(r)["score"]
	if !ok {
		http.Error(w, "invalid score", http.StatusBadRequest)
		return
	}
	pb, err := strconv.Atoi(progressBar)
	if err != nil {
		http.Error(w, "invalid score", http.StatusBadRequest)
		return
	}
	if pb > 100 {
		pb = 100
	}
	if pb < 0 {
		pb = 0
	}
	tmpl := template.Must(template.ParseFiles("template.html"))

	data := struct {
		Progress int
	}{
		Progress: pb,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
	}
}
