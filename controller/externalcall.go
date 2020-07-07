package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/reviashko/logopassapi/utils"
)

//ExternalLogic interface
type ExternalLogic interface {
	GetResult(io.ReadCloser) (string, error)
}

//ExternalCall struct
type ExternalCall struct {
	Cntrl         Controller
	ExternalLogic ExternalLogic
}

//CheckTokenAndDoFunc func
func (ec *ExternalCall) CheckTokenAndDoFunc(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		return
	}

	checked, err := ec.Cntrl.Crypto.CheckAuthToken(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Token validation error!",
			""))
		return
	}

	if !checked {
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Невалидный токен!",
			""))
		return
	}

	data, err := ec.ExternalLogic.GetResult(r.Body)
	if err != nil {
		fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
			false,
			"Ошибка обработки данных!",
			""))
		return
	}

	fmt.Fprintf(w, "%s", utils.GetJSONAnswer("",
		true,
		"",
		data))
}
