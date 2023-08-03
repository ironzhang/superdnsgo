package superutil

import (
	"encoding/json"
	"io/ioutil"
)

// ReadJSON read value from json file
func ReadJSON(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

// WriteJSON write value to json file
func WriteJSON(filename string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, data, 0666); err != nil {
		return err
	}
	return nil
}
