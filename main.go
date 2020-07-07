package main

import (
	"log"
	"net/http"

	"github.com/reviashko/logopassapi/controller"
)

func main() {

	cntrl := controller.NewController("config/db.json", "config/crypto.json", "config/smtp.json")
	log.Fatal(http.ListenAndServe(":8080", cntrl.GetRouter()))
}
