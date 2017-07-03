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

const Caroni23 = "caroni23"

var persistence map[string]domain.Rum

func init() {
	persistence = make(map[string]domain.Rum, 100)
	persistence[Caroni23] = domain.Rum {
		Name:"Caroni of Trinidad",
		BottlingDate: time.Date(1997, 02, 01, 0,0,0,0,time.UTC),
		Age:23,
	}
}


func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("entering root handler")
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	pusher, ok := w.(http.Pusher)
	if ok {
		pusher.Push("style.css", nil)
	}

	t := template.New("gopherTemplate")
	t, err := t.ParseFiles("template/gopher.tmpl")
	if err != nil {
		log.Printf("Error parsing template %v", err)
	}

	err = t.ExecuteTemplate(w, "gopher.tmpl", persistence[Caroni23])
	if err != nil {
		log.Printf("Error executing template %v", err)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("entering data handler")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	pathParams := strings.Split(r.URL.Path[1:], "/")
	log.Println(r.URL.Path)
	log.Printf("%+v", pathParams)

	/*
	queryParam, err := url.ParseQuery(r.URL.RawQuery)
	log.Printf("%+v", queryParam)
	*/

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
	log.Println("starting")

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/v1/rums/", dataHandler)
	http.HandleFunc("/index.html", rootHandler)

	err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
	//err := http.ListenAndServe(":8080", nil)
	if err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	log.Println("Server gracefully stopped")

	/*
	// subscribe to SIGINT signals
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	srv := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("could not shutdown: %v", err)
		}
	}()

	err := srv.ListenAndServeTLS("server.crt", "server.key")
	if err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	log.Println("Server gracefully stopped")
	*/
}
