package main

import (
	"log"
	"net/http"

	"github.com/reviashko/logopassapi/controller"
	"github.com/reviashko/logopassapi/example"
)

func main() {

	cntrl := controller.NewController("config/db.json", "config/crypto.json", "config/smtp.json")
	router := cntrl.NewRouter()

	//Example. Apply external logic after authtorization token check
	externalCallExample := controller.ExternalCall{Cntrl: cntrl, ExternalLogic: &example.ExternalLogicExample{}}
	router.HandleFunc("/gettestdatabytoken/", externalCallExample.CheckTokenAndDoFunc).Methods("POST", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", router))
}
