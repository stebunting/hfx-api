package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stebunting/hfx-backend/cors"
	"github.com/stebunting/hfx-backend/model"
	"github.com/stebunting/hfx-backend/routes"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")

	model := model.Model{}
	model.Connect()

	r := routes.ConfigRoutes(model.Db)

	http.Handle("/wake", cors.Middleware(http.HandlerFunc(r.Wake)))
	http.Handle("/dbinit", cors.Middleware(http.HandlerFunc(r.DbInit)))
	http.Handle("/updatecurrencies", cors.Middleware(http.HandlerFunc(r.UpdateCurrencies)))
	http.Handle("/getcurrencies", cors.Middleware(http.HandlerFunc(r.GetCurrencies)))
	http.Handle("/getrate", cors.Middleware(http.HandlerFunc(r.GetRate)))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
