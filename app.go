package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
}

type shortenReq struct {
	URL string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64 `json:"expiration_in_minutes" validate:"min=0"`
}

type shortlinkResp struct {
	Shortlink string `json:"shortlink"`
}

func (a *App) Initialize() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/shorten", a.createShortlink).Methods("POST")
	a.Router.HandleFunc("/api/info", a.getShortlinkInfo).Methods("GET")
	a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9][1,11]}", a.redirect).Methods("GET")
}

func (a *App) createShortlink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	if err := validator.New().Struct(req); err != nil {
		return
	}
	defer r.Body.Close()
	fmt.Printf("%v\n", req)
}

func (a *App) getShortlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")

	fmt.Printf("%s\n", s)
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s\n", vars["shortlink"])
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}