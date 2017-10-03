package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

const caroni = "caroni"
const worthyPark = "worthyPark"

type Rum struct {
	Name         string    `json:"name"`
	Age          int       `json:"age"`
	BottlingDate time.Time `json:"bottling_date"`
}

var persistence map[string]Rum
var tmpl *template.Template

func init() {
	persistence = make(map[string]Rum, 10)

	persistence[caroni] = Rum{
		Name:         "Caroni of Trinidad",
		BottlingDate: time.Date(1996, 02, 01, 0, 0, 0, 0, time.UTC),
		Age:          20,
	}

	persistence[worthyPark] = Rum{
		Name:         "Worthy Park of Jamaica",
		BottlingDate: time.Date(2007, 12, 30, 0, 0, 0, 0, time.UTC),
		Age:          9,
	}

	var err error
	tmpl = template.New("gopherTemplate")
	tmpl, err = tmpl.ParseFiles("template/gopher.tmpl.html")
	if err != nil {
		log.Printf("Error parsing template %v", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("entering root handler")

	pusher, ok := w.(http.Pusher)
	if ok {
		err := pusher.Push("/static/css/bootstrap.min.css", nil)
		if err != nil {
			log.Printf("Error pushing bootstrap css %v", err)
		}

		err = pusher.Push("/static/js/jquery-1.12.4.min.js", nil)
		if err != nil {
			log.Printf("Error pushing jquery %v", err)
		}

		err = pusher.Push("/static/js/bootstrap.min.js", nil)
		if err != nil {
			log.Printf("Error pushing bootstrap js %v", err)
		}

		err = pusher.Push("/static/img/gopher-dance-long.gif", nil)
		if err != nil {
			log.Printf("Error pushing gif %v", err)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := tmpl.ExecuteTemplate(w, "gopher.tmpl.html", persistence)
	if err != nil {
		log.Printf("Error executing template %v", err)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("entering data handler")

	bytes, err := json.Marshal(persistence[caroni])
	if err != nil {
		http.Error(w, "Error marshalling rum", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	w.Write(bytes)
}

func main() {
	log.Println("servers starting...")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/index.html", rootHandler)
	http.HandleFunc("/v1/rums/", dataHandler)

	log.Println("launching http on localhost:8081...")
	go func() {
		http.ListenAndServe(":8081", nil)
	}()

	log.Println("launching https on localhost:8080...")
	http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)

	log.Println("...servers stopped")
}
