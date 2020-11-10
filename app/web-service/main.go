package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/subosito/gotenv"

	"github.com/rbpermadi/whim_assignment/config"
	"github.com/rbpermadi/whim_assignment/handler"
)

func main() {
	gotenv.Load()

	db := config.NewMySQL()
	defer db.Close()

	h := handler.NewHandler()

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler:      h,
		ErrorLog:     logger,
	}

	log.Printf("whim is available at %s\n", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
