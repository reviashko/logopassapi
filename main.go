package main

/*
1. Добавить поддержку версии пакетов
2. Роутер подогнать под REST
3. Нужно бы еще middleware добавить с логированием
*/

import (
	"github.com/logopassapi/auth"
	"github.com/logopassapi/controllers"
	"github.com/logopassapi/models"
	"github.com/logopassapi/utils"

	"log"
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
