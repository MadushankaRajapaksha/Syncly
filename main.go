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

func roomRoute(w http.ResponseWriter, r *http.Request) {
	// get room id
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	templ, err := template.ParseFiles("web/room.html")
	if err != nil {
		fmt.Fprint(w, "ERORR : Internal Server ERORR")

	}
	templ.Execute(w, roomId)

}

func userRoomRoute(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("web/userRooms.html")
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
	// room route handler
	r.HandleFunc("/room/{roomId}", logging(roomRoute))
	// handler user rooms alredy used
	r.HandleFunc("/room", logging(userRoomRoute))

	// api routes
	r.HandleFunc("/api/room", logging(a.CreateRoom))
	r.HandleFunc("/api/room/{roomId}", logging(a.RoomHandler)).Methods("POST")
	r.HandleFunc("/sync/{roomId}", logging(api.SyncRoom))
	r.HandleFunc("/api/validate", logging(a.ValidateUser))
	r.HandleFunc("/api/user", logging(a.CreateUser)).Methods("POST")

	fmt.Printf("Server in runnig : http://localhost:%s \n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), r)

}
