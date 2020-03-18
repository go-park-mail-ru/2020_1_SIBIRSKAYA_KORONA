package message

import (
	"encoding/json"
)

type Pair struct {
	Name string
	Data interface{}
}

func GetBody(bodyData ...Pair) (string, error) {
	bodyMap := make(map[string]interface{})
	for _, elem := range bodyData {
		bodyMap[elem.Name] = elem.Data
	}
	bodyJSON, err := json.Marshal(bodyMap)
	return string(bodyJSON), err
}