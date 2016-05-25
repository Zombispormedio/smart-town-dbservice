package routes

import (
	"github.com/Zombispormedio/smartdb/consumer"
	"github.com/Zombispormedio/smartdb/controllers"
	"github.com/Zombispormedio/smartdb/lib/store"
)

func Consumer(deliver *consumer.Deliver) error {
	var Error error
	keyElem := deliver.KeyElem()

	var id string
	Error = store.Get("push_identifier", "Config", func(value string) {
		id = value
	})

	if Error != nil {
		return Error
	}

	switch keyElem.Parent {
	case id:
		Error = Push(deliver, keyElem.Children)

	}
	return Error

}

func Push(deliver *consumer.Deliver, children []string) error {
	var Error error

	result := map[string]interface{}{}
	BodyError := deliver.Body(&result)

	if BodyError != nil {
		return BodyError
	}

	session := deliver.Session()

	key := children[0]

	switch key {
	case "push":
		controllers.RabbitPushSensor(result, session)

	}

	return Error
}
