package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stebunting/hfx-backend/cors"
	"github.com/stebunting/hfx-backend/model"
	"github.com/stebunting/hfx-backend/server"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")

	model := model.Model{}
	model.Connect()

	server := server.ConfigRoutes(model.Db)

	http.Handle("/wake", cors.Middleware(http.HandlerFunc(server.Wake)))
	http.Handle("/dbinit", cors.Middleware(http.HandlerFunc(server.DbInit)))
	http.Handle("/updatecurrencies", cors.Middleware(http.HandlerFunc(server.UpdateCurrencies)))
	http.Handle("/getcurrencies", cors.Middleware(http.HandlerFunc(server.GetCurrencies)))
	http.Handle("/getrate", cors.Middleware(http.HandlerFunc(server.GetRate)))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
