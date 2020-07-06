package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
)

//GetJSONAnswer func
func GetJSONAnswer(token string, accepted bool, reason string, data string) string {
	return fmt.Sprintf(`{"accepted":%t, "token":"%s", "reason":"%s", "data":"%s"}`, accepted, token, reason, data)
}

//ConvertBody2JSON func
func ConvertBody2JSON(data io.Reader, v interface{}) error {
	body, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(string(body)), v)
	if err != nil {
		return err
	}

	return nil
}

//CheckEmailFormat func
func CheckEmailFormat(email string) bool {

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
}
