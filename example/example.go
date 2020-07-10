package example

import (
	"io"

	"github.com/reviashko/logopassapi/auth"

	"github.com/reviashko/logopassapi/utils"
)

//ExternalLogicExample struct
type ExternalLogicExample struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

//GetResult func
func (ele *ExternalLogicExample) GetResult(requestBody io.ReadCloser, token auth.Token, httpMethod string) (string, error) {
	var obj ExternalLogicExample
	err := utils.ConvertBody2JSON(requestBody, &obj)
	if err != nil {
		return "", err
	}

	if obj.Name != "test" || token.UserID != 1 {
		return `{"result":"error"}`, nil
	}

	return `{"result":"ok"}`, nil
}
