package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/reviashko/logopassapi/auth"
	"github.com/reviashko/logopassapi/models"
	"github.com/reviashko/logopassapi/utils"

	"github.com/gorilla/mux"
)

//Controller struct
type Controller struct {
	Db     models.Datastore
	Crypto auth.CryptoData
	SMTP   utils.SMTPData
}

//NewController func
func NewController(db models.Datastore, cryptoData auth.CryptoData, smtpData utils.SMTPData) Controller {
	return Controller{Db: db, Crypto: cryptoData, SMTP: smtpData}
}

//NewRouter func
func (c *Controller) NewRouter() *mux.Router {

	router := mux.NewRouter()

	//login method - tested
	router.HandleFunc("/getauthtoken/", c.GetAuthTokenHandler).Methods("POST")
	//registration method -tested
	router.HandleFunc("/registration/", c.RegistrationHandler).Methods("POST")
	//forgot password method (sending special restore link) -tested
	router.HandleFunc("/getpasswordrestoreemail/", c.SendRestorePasswordEmailHandler).Methods("POST")
	//change password method -tested
	router.HandleFunc("/changepassword/{token}", c.ChangePasswordHandler).Methods("GET")

	return router
}

//ChangePasswordHandler func
func (c *Controller) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	linkData := vars["token"]

	if len(linkData) == 0 {
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Empty token!",
			""))
		return
	}

	tokenJSON, err := c.Crypto.DecryptTextAES256(linkData)
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Token decryption error!",
			""))
		return
	}

	var token auth.Token
	err = json.Unmarshal([]byte(tokenJSON), &token)
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Некорректная ссылка!",
			""))
		return
	}

	if token.TTL-time.Now().Unix() > 0 {
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Просроченная ссылка!",
			""))
		return
	}

	userData, _, err := c.Db.GetUserByEmail(token.Email)
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Некорректный email!",
			""))
		return
	}

	password, err := auth.GetNewPassword()
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Password generation error!",
			""))
		return
	}

	userData.PswdHashB = c.Crypto.GetSHA256Bytes(password)

	userID, _, err := c.Db.SaveUser(userData)
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Save new password error!",
			""))
		return
	}

	if userID > 0 {

		err = c.SMTP.SendEmail(userData.Email, `Subject: Ваш пароль\n `+password)
		if err != nil {

			log.Println(err.Error())
			//fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			//	false,
			//	err.Error(),
			//	""))
			//return
			//TODO: some logs here
		}

		token, _ := c.Crypto.EncryptTextAES256Base64(c.Crypto.GetTokenJSON(userID))

		fmt.Fprintf(w, "%s", utils.GetJSONAnswer(token,
			true,
			"",
			""))

	} else {
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Ошибка смены пароля!",
			""))
	}

}

//SendRestorePasswordEmailHandler func
func (c *Controller) SendRestorePasswordEmailHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var tokenItem auth.Token
	err := utils.ConvertBody2JSON(r.Body, &tokenItem)
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Wrong data!",
			""))
		return
	}

	//check email format
	if !utils.CheckEmailFormat(tokenItem.Email) {
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Некорректный EMail формат!",
			""))
		return
	}

	linkData, err := c.Crypto.EncryptTextAES256Base64(fmt.Sprintf(`{"email":"%s", "ttl":%d}`, tokenItem.Email, c.Crypto.PasswordEmailTTL))
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Decrypt error!",
			""))
		return
	}

	err = c.SMTP.SendEmail(tokenItem.Email, `Subject: Смена пароля: `+c.Crypto.RestorePasswordURL+linkData)
	if err != nil {

		log.Println(err.Error())
		//fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
		//	false,
		//	"EMail link error!",
		//	""))
		//return
		//TODO: some logs here
	}

	fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
		true,
		"Вам отправлен EMail с инструкциями!",
		""))
}

//RegistrationHandler func
func (c *Controller) RegistrationHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	userData := new(models.UserData)
	err := utils.ConvertBody2JSON(r.Body, &userData)
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Bad data",
			""))
		return
	}

	userData.Email = strings.ToLower(userData.Email)

	//check email format
	if !utils.CheckEmailFormat(userData.Email) {
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Некорректный EMail формат!",
			""))
		return
	}

	//generate password and send it to email
	password, err := auth.GetNewPassword()
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Password error",
			""))
		return
	}

	userData.PswdHashB = c.Crypto.GetSHA256Bytes(password)

	userID, errorCode, err := c.Db.SaveUser(userData)
	if err != nil && errorCode != "22024" {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Save error!",
			""))
		return
	}

	if userID > 0 {

		err = c.SMTP.SendEmail(userData.Email, `Subject: Ваш пароль\n`+password)
		if err != nil {

			log.Println(err.Error())
			//fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			//	false,
			//	err.Error(),
			//	""))
			//return
			//TODO: logs here
		}

		token, _ := c.Crypto.EncryptTextAES256Base64(c.Crypto.GetTokenJSON(userID))

		fmt.Fprintf(w, "%s", utils.GetJSONAnswer(token,
			true,
			"",
			""))
	} else {
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Такой EMail уже используется!",
			""))
	}
}

//GetAuthTokenHandler func
func (c *Controller) GetAuthTokenHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var afd auth.LogoPassData
	err := utils.ConvertBody2JSON(r.Body, &afd)
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Wrong data!",
			""))
		return
	}

	userData, _, err := c.Db.GetUserByAuth(afd.Login, c.Crypto.GetSHA256Bytes(afd.Password))
	if err != nil {

		log.Println(err.Error())
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"User not found!",
			""))
		return
	}

	if userData.UserID > 0 {

		token, _ := c.Crypto.EncryptTextAES256Base64(c.Crypto.GetTokenJSON(userData.UserID))

		fmt.Fprintf(w, "%s", utils.GetJSONAnswer(token,
			true,
			"",
			""))

	} else {

		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Не верный логин или пароль!",
			""))
	}

}
