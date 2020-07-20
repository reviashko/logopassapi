package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
)

//SortedIntInsert func
func SortedIntInsert(arr []int, val int, prevSearchIndex int) []int {
	//i := sort.SearchInts(arr, val)
	arr = append(arr, 0)
	copy(arr[prevSearchIndex+1:], arr[prevSearchIndex:])
	arr[prevSearchIndex] = val
	return arr
}

//SortedStringInsert func
func SortedStringInsert(arr []string, val string, prevSearchIndex int) []string {
	//i := sort.SearchStrings(arr, val)
	arr = append(arr, "")
	copy(arr[prevSearchIndex+1:], arr[prevSearchIndex:])
	arr[prevSearchIndex] = val
	return arr
}

//GetJSONAnswer func
func GetJSONAnswer(token string, accepted bool, reason string, jsonData string) string {
	qv := `"`
	if len(jsonData) > 0 {
		qv = ``
	}
	return fmt.Sprintf(`{"accepted":%t, "token":"%s", "reason":"%s", "data":%s%s%s}`, accepted, token, reason, qv, jsonData, qv)
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
