package main

import "encoding/json"

func jsonToMap(dataByte []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(dataByte, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
