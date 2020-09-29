package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

//CryptoData struct
type CryptoData struct {
	AES256Key          string `json:"AES256Key"`
	SHA256Salt         string `json:"SHA256Salt"`
	TokenTTL           int64  `json:"TokenTTL"`
	PasswordEmailTTL   int64  `json:"PasswordEmailTTL"`
	RestorePasswordURL string `json:"RestorePasswordURL"`
	CaptchaTTL         int64  `json:"CaptchaTTL"`
}

//GetTokenJSON func
func (c *CryptoData) GetTokenJSON(userID int) string {
	return fmt.Sprintf(`{"user_id":%d, "ttl":%d}`, userID, time.Now().Unix()+c.TokenTTL)
}

//CheckConfig func
func (c *CryptoData) CheckConfig() error {

	if len(c.AES256Key) != 32 {
		return errors.New("Crypto key is not a multiple of 32")
	}

	return nil
}

//CheckAuthToken func
func (c *CryptoData) CheckAuthToken(authHeaderValue string) (bool, Token, error) {

	tokenJSON, _ := c.DecryptTextAES256(strings.ReplaceAll(authHeaderValue, `"`, ""))

	var token Token
	err := json.Unmarshal([]byte(tokenJSON), &token)
	if err != nil {

		log.Println(err.Error())
		return false, token, err
	}

	return (token.TTL-time.Now().Unix() > 0), token, err
}

//GetSHA256Bytes func
func (c *CryptoData) GetSHA256Bytes(text string) []byte {
	h := sha256.New()
	h.Write([]byte(text + c.SHA256Salt))
	return h.Sum(nil)
}

//EncryptTextAES256Base64 func
func (c *CryptoData) EncryptTextAES256Base64(textString string) (string, error) {

	if len(c.AES256Key) != 32 {
		panic("too short key!")
	}

	text := []byte(textString)
	key := []byte(c.AES256Key)

	cp, err := aes.NewCipher(key)
	if err != nil {

		log.Println(err.Error())
		return "", err
	}

	gcm, err := cipher.NewGCM(cp)
	if err != nil {

		log.Println(err.Error())
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {

		log.Println(err.Error())
		return "", err
	}

	txt := gcm.Seal(nonce, nonce, text, nil)

	return b64.RawURLEncoding.EncodeToString([]byte(txt)), nil
}

//DecryptTextAES256 func
func (c *CryptoData) DecryptTextAES256(encryptedBase64 string) (string, error) {

	key := []byte(c.AES256Key)

	ciphertext, err := b64.RawURLEncoding.DecodeString(encryptedBase64) //[]byte(encryptedText)
	if err != nil {

		log.Println(err.Error())
		return "", err
	}

	cp, err := aes.NewCipher(key)
	if err != nil {

		log.Println(err.Error())
		return "", err
	}

	gcm, err := cipher.NewGCM(cp)
	if err != nil {

		log.Println(err.Error())
		return "", nil
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {

		log.Println(err.Error())
		return "", nil
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {

		log.Println(err.Error())
		return "", nil
	}

	return string(plaintext), nil
}
