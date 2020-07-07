package main

import (
	"log"
	"net/http"

	"github.com/reviashko/logopassapi/controller"
)

func main() {

	controller := controller.NewController("config/db.json", "config/crypto.json", "config/smtp.json")
	log.Fatal(http.ListenAndServe(":8080", controller.GetRouter()))
}
