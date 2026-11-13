package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"syncly/api"
	"syncly/core"

	"github.com/gorilla/mux"
)

const port string = "3000"

// logging midelware basc setup

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}

// ###############################
// "/" route setup
// ###############################
func landingPage(w http.ResponseWriter, r *http.Request) {
	// loder html file from web folder and load as template

	templ, err := template.ParseFiles("web/index.html")
	// handler err from laod file
	if err != nil {
		fmt.Fprint(w, "ERORR : Load landing page!")

	}

	// test parese backend data to html file
	fromApp := "This fro backend"

	templ.Execute(w, fromApp)

}

func createRoomRoute(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("web/create-room.html")
	// handler err from laod file
	if err != nil {
		fmt.Fprint(w, "ERORR : Internal Server ERORR")

	}

	templ.Execute(w, "")

}

func main() {

	// connect to database

	pool := core.ConnectDb()
	a := api.New(pool)

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("assets/"))

	// static files router
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// defined landigp page to sevre
	r.HandleFunc("/", logging(landingPage))
	// definde createroom route
	r.HandleFunc("/create-room", logging(createRoomRoute))

	// api routes
	r.HandleFunc("/api/room", logging(a.CreateRoom))

	fmt.Printf("Server in runnig : http://localhost:%s \n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), r)

}
