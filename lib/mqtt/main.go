package mqtt

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

func CreateUser(user map[string]interface{}) error {
	var Error error

	CreateError := Request("POST", "user", user)

	if CreateError != nil {
		return CreateError
	}

	acl := map[string]interface{}{}

	acl["username"] = user["username"]

	acl["read"] = false
	acl["write"] = true
	acl["topic"] = os.Getenv("MQTT_TOPIC")

	ACLError := Request("POST", "acl", acl)

	if ACLError != nil {
		return ACLError
	}

	return Error

}

func DelUser(id string) error {
	var Error error

	DelError := Request("DELETE", "user/"+id, nil)

	if DelError != nil {
		return DelError
	}

	acl := map[string]interface{}{}

	acl["username"] = id
	acl["topic"] = os.Getenv("MQTT_TOPIC")

	ACLError := Request("DELETE", "acl", acl)

	if ACLError != nil {
		return ACLError
	}

	return Error
}

func Request(method string, param string, data map[string]interface{}) error {
	var Error error

	url := os.Getenv("MQTT_API") + param

	var req *http.Request

	if data != nil {
		body, JSONError := json.Marshal(data)

		if JSONError != nil {
			return  JSONError
		}

		payload := strings.NewReader(string(body))
		req, _ = http.NewRequest(method, url, payload)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	req.Header.Add("authorization", os.Getenv("MQTT_AUTH"))
	req.Header.Add("content-type", "application/json")

	_, Error = http.DefaultClient.Do(req)

	return Error
}
