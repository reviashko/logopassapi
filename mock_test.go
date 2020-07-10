package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/reviashko/logopassapi/auth"
	"github.com/reviashko/logopassapi/controller"
	"github.com/reviashko/logopassapi/example"
	"github.com/reviashko/logopassapi/models"
	"github.com/reviashko/logopassapi/utils"

	"github.com/gorilla/mux"

	"github.com/lib/pq"
	"github.com/tkanos/gonfig"
)

type mockDB struct{}

func (mdb *mockDB) GetUser(userID int) (*models.UserData, pq.ErrorCode, error) {

	smtpData := utils.SMTPData{}
	if gonfig.GetConf("config/smtp.json", &smtpData) != nil {
		log.Panic("load smtp confg error")
	}

	userData := new(models.UserData)
	userData.UserID = 1
	userData.IsActive = true
	userData.FirstName = "test"
	userData.LastName = "test"
	userData.Email = smtpData.MockEmail

	var errorCode pq.ErrorCode

	return userData, errorCode, nil
}

func (mdb *mockDB) GetUserByEmail(email string) (*models.UserData, pq.ErrorCode, error) {
	return mdb.GetUser(1)
}

func (mdb *mockDB) GetUserByAuth(email string, pswdHashB []byte) (*models.UserData, pq.ErrorCode, error) {

	return mdb.GetUser(1)
}

func (mdb *mockDB) SaveUser(userData *models.UserData) (int, pq.ErrorCode, error) {

	var errorCode pq.ErrorCode

	return 1, errorCode, nil
}

func InitController(mdb *mockDB) controller.Controller {

	cryptoData := auth.CryptoData{}
	if gonfig.GetConf("config/crypto.json", &cryptoData) != nil {
		log.Panic("load crypto confg error")
	}

	smtpData := utils.SMTPData{}
	if gonfig.GetConf("config/smtp.json", &smtpData) != nil {
		log.Panic("load smtp confg error")
	}

	return controller.Controller{Db: mdb, Crypto: cryptoData, SMTP: smtpData}
}

//TestGetTestDataByTokenHandler function
func TestGetTestDataByTokenHandler(t *testing.T) {

	cntrl := InitController(&mockDB{})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/gettestdatabytoken/", bytes.NewReader([]byte(`{"name": "test", "text": "test"}`)))

	token, _ := cntrl.Crypto.EncryptTextAES256Base64(cntrl.Crypto.GetTokenJSON(1))
	req.Header.Add("Authorization", token)

	externalCallExample := controller.ExternalCall{Cntrl: cntrl, ExternalLogic: &example.ExternalLogicExample{}}
	http.HandlerFunc(externalCallExample.CheckTokenAndDoFunc).ServeHTTP(rec, req)

	expected := `{"accepted":true, "token":"", "reason":"", "data":{"result":"ok"}}`
	if expected != rec.Body.String() {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	}
}

//TestGetAuthTokenHandler func
func TestGetAuthTokenHandler(t *testing.T) {

	controller := InitController(&mockDB{})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/getauthtoken/", bytes.NewReader([]byte(`{"login": "`+controller.SMTP.MockEmail+`", "password": "test"}`)))

	http.HandlerFunc(controller.GetAuthTokenHandler).ServeHTTP(rec, req)

	var answer auth.JSONAnswer
	err := utils.ConvertBody2JSON(rec.Body, &answer)
	if err != nil {
		t.Errorf("\n...ConvertBody2JSON error = %s\n", err.Error())
	}

	if !answer.Accepted {
		t.Errorf("\n...Wrong Accepted answer parametr error")
	} else {

		decryptedToken, err := controller.Crypto.DecryptTextAES256(answer.Token)
		if err != nil {
			t.Errorf("\n...Token decryption error = %s\n", err.Error())
		}

		var token auth.Token
		err = json.Unmarshal([]byte(decryptedToken), &token)
		if err != nil {
			t.Errorf("\n...Token Unmarshal error = %s\n", err.Error())
		}

		if token.TTL-time.Now().Unix() < 0 || token.TTL-time.Now().Unix() > controller.Crypto.TokenTTL {
			t.Errorf("\n...Wrong token ttl error")
		}

		if token.UserID != 1 {
			t.Errorf("\n...Wrong token user_id error = %s\n", err.Error())
		}
	}

}

//TestRegistrationHandler func
func TestRegistrationHandler(t *testing.T) {

	controller := InitController(&mockDB{})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/registration/", bytes.NewReader([]byte(`{"email": "`+controller.SMTP.MockEmail+`", "first_name": "test", "is_active": true}`)))

	http.HandlerFunc(controller.RegistrationHandler).ServeHTTP(rec, req)

	var answer auth.JSONAnswer
	err := utils.ConvertBody2JSON(rec.Body, &answer)
	if err != nil {
		t.Errorf("\n...ConvertBody2JSON error = %s\n...Body=%v", err.Error(), rec.Body.String())
	}

	if answer.Accepted {
		decryptedToken, err := controller.Crypto.DecryptTextAES256(answer.Token)
		if err != nil {
			t.Errorf("\n...Token decryption error = %s\n", err.Error())
		}

		var token auth.Token
		err = json.Unmarshal([]byte(decryptedToken), &token)
		if err != nil {
			t.Errorf("\n...Token Unmarshal error = %s\n...Reason=%v", err.Error(), answer.Reason)
		}

		if token.TTL-time.Now().Unix() < 0 || token.TTL-time.Now().Unix() > controller.Crypto.TokenTTL {
			t.Errorf("\n...Wrong token ttl error")
		}

		if token.UserID != 1 {
			t.Errorf("\n...Wrong token user_id error = %s\n", err.Error())
		}
	} else {
		t.Errorf("\n...Wrong Accepted answer parametr error")
	}

}

//TestSendRestorePasswordEmailHandler func
func TestSendRestorePasswordEmailHandler(t *testing.T) {

	controller := InitController(&mockDB{})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/getpasswordrestoreemail/", bytes.NewReader([]byte(`{"email": "`+controller.SMTP.MockEmail+`"}`)))

	http.HandlerFunc(controller.SendRestorePasswordEmailHandler).ServeHTTP(rec, req)

	var answer auth.JSONAnswer
	err := utils.ConvertBody2JSON(rec.Body, &answer)
	if err != nil {
		t.Errorf("\n...ConvertBody2JSON error = %s\n", err.Error())
	}

	if !answer.Accepted {
		t.Errorf("\n...Wrong Accepted parametr error")
	}

}

//TestSendRestorePasswordEmailHandler func
func TestChangePasswordHandler(t *testing.T) {

	controller := InitController(&mockDB{})

	linkData, err := controller.Crypto.EncryptTextAES256Base64(fmt.Sprintf(`{"email":"%s", "ttl":%d}`, controller.SMTP.MockEmail, controller.Crypto.PasswordEmailTTL))
	if err != nil {
		t.Errorf("\n...Restore password link generation error = %s\n", err.Error())
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/changepassword/", nil)

	req = mux.SetURLVars(req, map[string]string{
		"token": linkData,
	})

	http.HandlerFunc(controller.ChangePasswordHandler).ServeHTTP(rec, req)

	var answer auth.JSONAnswer
	err = utils.ConvertBody2JSON(rec.Body, &answer)
	if err != nil {
		t.Errorf("\n...ConvertBody2JSON error = %s\n", err.Error())
	}

	if !answer.Accepted {
		t.Errorf("\n...Wrong Accepted parametr error\n...Reason=%v\n...Link=%v", answer.Reason, linkData)
	} else {
		decryptedToken, err := controller.Crypto.DecryptTextAES256(answer.Token)
		if err != nil {
			t.Errorf("\n...Token decryption error = %s\n", err.Error())
		}

		var token auth.Token
		err = json.Unmarshal([]byte(decryptedToken), &token)
		if err != nil {
			t.Errorf("\n...Token Unmarshal error = %s\n...Reason=%v", err.Error(), answer.Reason)
		}

		if token.TTL-time.Now().Unix() < 0 || token.TTL-time.Now().Unix() > controller.Crypto.TokenTTL {
			t.Errorf("\n...Wrong token ttl error")
		}

		if token.UserID != 1 {
			t.Errorf("\n...Wrong token user_id error = %s\n", err.Error())
		}
	}

}
