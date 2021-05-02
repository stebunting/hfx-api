package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stebunting/hfx-backend/model"
	"github.com/stebunting/hfx-backend/routes"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")

	model := model.Model{}
	model.Connect()

	r := routes.ConfigRoutes(model.Db)

	http.HandleFunc("/wake", r.Wake)
	http.HandleFunc("/dbinit", r.DbInit)
	http.HandleFunc("/updatecurrencies", r.UpdateCurrencies)
	http.HandleFunc("/getcurrencies", r.GetCurrencies)
	http.HandleFunc("/getrate", r.GetRate)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
