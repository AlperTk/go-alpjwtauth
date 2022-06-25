package controllers

import (
	"encoding/json"
	"github.com/AlperTk/go-jwt-role-based-auth/src/example/authentication/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type EventController struct{}

var prefixUrl = "/api/v1"

func (t EventController) RegisterRoutes(router *mux.Router) {
	router.Handle(prefixUrl+"/test", http.HandlerFunc(postTest)).Methods("POST")
}

func postTest(w http.ResponseWriter, _ *http.Request) {
	response, err := services.CreateTest()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
