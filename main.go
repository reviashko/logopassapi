package main

import (
	"log"
	"logopassapi/controllers"
	"net/http"
)

func main() {

	controller := controllers.Init("config/db.json", "config/crypto.json", "config/smtp.json")
	log.Fatal(http.ListenAndServe(":8080", controller.GetRouter()))
}
