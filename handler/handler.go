package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rbpermadi/whim_assignment/app/response"
	"github.com/rs/cors"
)

type Registration interface {
	Register(r *httprouter.Router) error
}

func Healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "ok")
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	meta := response.MetaInfo{HTTPStatus: 404}
	res := response.ResponseBody{Message: "path not found", Meta: meta}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(meta.HTTPStatus)
	json.NewEncoder(w).Encode(res)
}

func NewHandler(registrations ...Registration) http.Handler {
	router := httprouter.New()
	router.HandleMethodNotAllowed = false

	router.HandlerFunc("GET", "/healthz", Healthz)
	// start route
	for _, reg := range registrations {
		reg.Register(router)
	}

	router.NotFound = http.HandlerFunc(NotFound)

	co := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		MaxAge:         86400,
	})

	return co.Handler(router)
}
