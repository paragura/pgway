package pgway

import (
	"encoding/json"
)

func CreateJsonString(i interface{}) (string, error) {

	if i == nil {
		return "", nil
	}

	buf, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
