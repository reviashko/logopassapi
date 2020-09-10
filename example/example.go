package example

import (
	"errors"
	"log"
	"net/http"

	"github.com/reviashko/logopassapi/auth"

	"github.com/reviashko/logopassapi/utils"
)

//ExternalLogicExample struct
type ExternalLogicExample struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

//GetResult func
func (ele *ExternalLogicExample) GetResult(request *http.Request, token auth.Token) (string, error) {
	var obj ExternalLogicExample
	err := utils.ConvertBody2JSON(request.Body, &obj)
	if err != nil {

		log.Println(err.Error())
		return "", err
	}

	if obj.Name != "test" || token.UserID != 1 {
		return `{"result":"error"}`, errors.New("wrong user")
	}

	return `{"result":"ok"}`, nil
}
