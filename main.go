package main

import (
	"github.com/gorilla/mux"
	"github.com/shitpostingio/analysis-api/api/handler"
	"github.com/shitpostingio/analysis-commons/health-check"
	"log"
	"net/http"
)

const (
	cfgPathKey     = "API_CFG_PATH"
	bindAddressKey = "API_BIND_ADDRESS"
)

var (
	cfgPath         = "config.toml"
	bindAddress     = "localhost:9999"
	telegramHandler *handler.Handler
	directHandler   *handler.Handler
	r               *mux.Router
)

func main() {

	r.HandleFunc("/analysis/{reqType}/{type}/{id}", telegramHandler.Dispatch).Methods("GET")
	r.HandleFunc("/analysis/{reqType}", telegramHandler.Dispatch).Methods("GET")
	r.HandleFunc("/analysis/upload/{reqType}/{type}/{id}", directHandler.Dispatch).Methods("POST")
	r.HandleFunc("/analysis/upload/{reqType}", directHandler.Dispatch).Methods("POST")
	r.HandleFunc("/healthy", health_check.ConfirmServiceHealth).Methods("GET")

	log.Println("Analysis API serving")
	err := http.ListenAndServe(bindAddress, r)
	if err != nil {
		log.Fatal("Unable to listen and serve", err)
	}

}
