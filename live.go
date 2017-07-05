package main

import (
	"time"
	"net/http"
	"encoding/json"
	"log"
	"github.com/sebastienfr/golanglive/domain"
	"strings"
	"html/template"
)

const Caroni = "caroni"
const WorthyPark = "worthyPark"

var persistence map[string]domain.Rum

func init() {
	persistence = make(map[string]domain.Rum, 100)

	persistence[Caroni] = domain.Rum {
		Name:"Caroni of Trinidad",
		BottlingDate: time.Date(1996, 02, 01, 0,0,0,0,time.UTC),
		Age:20,
	}

	persistence[WorthyPark] = domain.Rum {
		Name:"Worthy Park of Jamaica",
		BottlingDate: time.Date(2007, 12, 30, 0,0,0,0,time.UTC),
		Age:9,
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("entering root handler")
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	pusher, ok := w.(http.Pusher)
	if ok {
		err := pusher.Push("/static/css/bootstrap.min.css", nil)
		if err != nil {
			log.Printf("Error pushing bootstrqp css %v", err)
		}

		err = pusher.Push("/static/js/jquery-1.12.4.min.js", nil)
		if err != nil {
			log.Printf("Error pushing jquery %v", err)
		}

		err =pusher.Push("/static/js/bootstrap.min.js", nil)
		if err != nil {
			log.Printf("Error pushing bootstrap js %v", err)
		}

		err =pusher.Push("/static/img/gopher-dance-long.gif", nil)
		if err != nil {
			log.Printf("Error pushing gif %v", err)
		}
	}

	t := template.New("gopherTemplate")
	t, err := t.ParseFiles("template/gopher.tmpl.html")
	if err != nil {
		log.Printf("Error parsing template %v", err)
	}

	err = t.ExecuteTemplate(w, "gopher.tmpl.html", persistence)
	if err != nil {
		log.Printf("Error executing template %v", err)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("entering data handler")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	pathParams := strings.Split(r.URL.Path[1:], "/")

	rum, ok := persistence[pathParams[2]]

	if !ok {
		http.Error(w, "Error rum not found", http.StatusNotFound)
		return
	}

	bytes, err  := json.Marshal(rum)
	if err != nil {
		http.Error(w, "Error marshalling rum", http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}


func main() {
	log.Println("servers starting...")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/v1/rums/", dataHandler)
	http.HandleFunc("/index.html", rootHandler)

	go func() {
		http.ListenAndServe(":8081", nil)
	}()
	log.Println("http ready on localhost:8081...")

	http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
	log.Println("https ready on localhost:8080...")

	log.Println("...servers stopped")
}
