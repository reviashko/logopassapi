package main

import (
	"log"
	"net/http"

	"github.com/reviashko/logopassapi/auth"
	"github.com/reviashko/logopassapi/controller"
	"github.com/reviashko/logopassapi/example"
	"github.com/reviashko/logopassapi/models"
	"github.com/reviashko/logopassapi/utils"
	"github.com/tkanos/gonfig"
)

func main() {

	connectionData := models.ConnectionData{}
	if gonfig.GetConf("config/db.json", &connectionData) != nil {
		log.Panic("load db confg error")
	}

	db, err := models.InitDB(connectionData)
	if err != nil {
		log.Panic(err)
	}

	cryptoData := auth.CryptoData{}
	if gonfig.GetConf("config/crypto.json", &cryptoData) != nil {
		log.Panic("load crypto confg error")
	}
	err = cryptoData.CheckConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	smtpData := utils.SMTPData{}
	if gonfig.GetConf("config/smtp.json", &smtpData) != nil {
		log.Panic("load smtp confg error")
	}

	frontSettings := auth.FrontSettings{}
	if gonfig.GetConf("config/front.json", &frontSettings) != nil {
		log.Panic("load front confg error")
	}

	cntrl := controller.NewController(db, cryptoData, smtpData, frontSettings)
	router := cntrl.NewRouter()

	//Example. Apply external logic after authtorization token checking
	externalCallExample := controller.ExternalCall{Cntrl: cntrl, ExternalLogic: &example.ExternalLogicExample{}, FrontSettings: frontSettings}
	router.HandleFunc("/gettestdatabytoken/", externalCallExample.CheckTokenAndDoFunc).Methods("POST", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", router))
}
