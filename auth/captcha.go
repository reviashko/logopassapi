package auth

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

//CAPTCHAData struct
type CAPTCHAData struct {
	Captcha    string `json:"captcha"`
	CaptchaURL string `json:"captcha_url"`
}

//Check func
func (c *CAPTCHAData) Check(decryptedCaptchaURL string) error {

	ttl, captchaValue, err := c.Parse(decryptedCaptchaURL)
	if err != nil {
		return errors.New("ОШИБКА В ДАННЫХ CAPTCHA")
	}

	if ttl-time.Now().Unix() < 0 {
		return errors.New("ИСТЕК СРОК ГОДНОСТИ CAPTCHA")
	}

	if len(c.Captcha) == 0 || c.Captcha != captchaValue {
		return errors.New("НЕВЕРНЫЙ ТЕКСТ КАРТИНКИ")
	}

	return nil
}

//Parse func
func (c *CAPTCHAData) Parse(decryptedCaptchaURL string) (int64, string, error) {

	tmpCaptchaArray := strings.SplitN(decryptedCaptchaURL, "#", 2)
	ttl, err := strconv.ParseInt(tmpCaptchaArray[1], 10, 64)
	if err != nil {
		return 0, "", errors.New("ОШИБКА В ДАННЫХ CAPTCHA")
	}

	return ttl, tmpCaptchaArray[0], nil
}
