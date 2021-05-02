package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/stebunting/hfx-backend/model"
	"github.com/stebunting/hfx-backend/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	model := model.Model{}
	model.Connect()

	r := routes.NewRoutes(model.Db)

	http.HandleFunc("/getrate", r.GetRate)
	http.HandleFunc("/getcurrencies", r.GetCurrencies)
	http.HandleFunc("/updatecurrencies", r.UpdateCurrencies)
	http.HandleFunc("/dbinit", r.DbInit)
	http.ListenAndServe(":5000", nil)
}
