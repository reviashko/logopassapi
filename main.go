package main

import (
	"log"
	"logopassapi/auth"
	"logopassapi/controllers"
	"logopassapi/models"
	"logopassapi/utils"
	"net/http"

	"github.com/tkanos/gonfig"
)

func main() {

	connectionData := models.ConnectionData{}
	if gonfig.GetConf("config/db.json", &connectionData) != nil {
		log.Panic("load db confg error")
	}

	cryptoData := auth.CryptoData{}
	if gonfig.GetConf("config/crypto.json", &cryptoData) != nil {
		log.Panic("load crypto confg error")
	}

	smtpData := utils.SMTPData{}
	if gonfig.GetConf("config/smtp.json", &smtpData) != nil {
		log.Panic("load smtp confg error")
	}

	db, err := models.InitDB(connectionData.ToString())
	if err != nil {
		log.Panic(err)
	}

	controller := controllers.Controllers{Db: db, Crypto: cryptoData, SMTP: smtpData}
	log.Fatal(http.ListenAndServe(":8080", controller.GetRouter()))
}
