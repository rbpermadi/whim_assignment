package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/subosito/gotenv"

	"github.com/rbpermadi/whim_assignment/config"
	"github.com/rbpermadi/whim_assignment/delivery"
	"github.com/rbpermadi/whim_assignment/handler"
	"github.com/rbpermadi/whim_assignment/repository"
	"github.com/rbpermadi/whim_assignment/usecase/currency"
)

func main() {
	gotenv.Load()

	db := config.NewMySQL()
	defer db.Close()

	// currencies
	currencyRepo := repository.NewMysqlCurrency(db)

	currencyUseCase := currency.NewService(&currency.Provider{
		Repo: currencyRepo,
	})

	currencyHandler := delivery.NewCurrencyHandler(currencyUseCase)

	h := handler.NewHandler(&currencyHandler)

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
