package main

import (
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

	r := routes.NewRoutes(model.Db)

	http.HandleFunc("/getrate", r.GetRate)
	http.HandleFunc("/getcurrencies", r.GetCurrencies)
	http.HandleFunc("/updatecurrencies", r.UpdateCurrencies)
	http.HandleFunc("/dbinit", r.DbInit)
	http.ListenAndServe(":"+port, nil)
}
