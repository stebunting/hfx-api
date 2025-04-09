package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stebunting/hfx-backend/cors"
	"github.com/stebunting/hfx-backend/model"
	"github.com/stebunting/hfx-backend/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}

	portEnv := os.Getenv("PORT")
	addressEnv := os.Getenv("DATABASE_URL")
	certFileEnv := os.Getenv("CERT_FILE")
	keyFileEnv := os.Getenv("KEY_FILE")

	port := flag.String("port", portEnv, "Port Number")
	dbPath := flag.String("db", addressEnv, "DB Path")
	fullchain := flag.String("fullchain", certFileEnv, "Full Chain Key")
	privKey := flag.String("privkey", keyFileEnv, "Private Key")

	flag.Parse()

	model := model.Model{}
	model.Connect(*dbPath)

	server := server.ConfigRoutes(model.Db)

	http.Handle("/wake", cors.Middleware(http.HandlerFunc(server.Wake)))
	http.Handle("/dbinit", cors.Middleware(http.HandlerFunc(server.DbInit)))
	http.Handle("/updatecurrencies", cors.Middleware(http.HandlerFunc(server.UpdateCurrencies)))
	http.Handle("/getcurrencies", cors.Middleware(http.HandlerFunc(server.GetCurrencies)))
	http.Handle("/getrate", cors.Middleware(http.HandlerFunc(server.GetRate)))

	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", *port), *fullchain, *privKey, nil))
}
