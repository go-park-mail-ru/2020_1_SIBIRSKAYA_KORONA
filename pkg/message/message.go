package message

import (
	"encoding/json"
)

type Pair struct {
	Name string
	Data interface{}
}

func GetBody(status uint, bodyData ...Pair) (string, error) {
	msg := make(map[string]interface{})
	msg["status"] = status
	if len(bodyData) != 0 {
		bodyMap := make(map[string]interface{})
		for _, elem := range bodyData {
			bodyMap[elem.Name] = elem.Data
		}
		msg["body"] = bodyMap
	}
	res, err := json.Marshal(msg)
	return string(res), err
}