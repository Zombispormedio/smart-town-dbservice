package rabbit

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Zombispormedio/smart-push/lib/utils"
	"github.com/streadway/amqp"
)

var EXCHANGE_TYPES = []interface{}{"headers", "direct", "topic", "fanout"}

type Rabbit struct {
	Conn *amqp.Connection

	Chan *amqp.Channel

	Queue *amqp.Queue

	ExName string
	ExType string
}

func New(exname string, extype string, durable bool) (*Rabbit, error) {
	var Error error
	rabbit := Rabbit{}
	rabbit.Conn, Error = amqp.Dial(os.Getenv("RABBIT"))

	if Error != nil {
		return nil, Error
	}

	rabbit.Chan, Error = rabbit.Conn.Channel()

	if Error != nil {
		return nil, Error
	}

	Error = rabbit.Exchange(exname, extype, durable)

	return &rabbit, Error
}

func (rabbit *Rabbit) Exchange(exname string, extype string, durable bool) error {
	var Error error

	if exname == "" {
		return errors.New("You Need Exname")
	}
	rabbit.ExName = exname

	if extype == "" {
		return errors.New("You Need ExType")
	}

	if !utils.Contains(EXCHANGE_TYPES, extype) {
		return errors.New("No valid Extype")
	}

	rabbit.ExType = extype

	Error = rabbit.Chan.ExchangeDeclare(
		rabbit.ExName, // name
		rabbit.ExType, // type
		durable,       // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	return Error

}

func (rabbit *Rabbit) Close() {
	rabbit.Chan.Close()
	rabbit.Conn.Close()
}

func (rabbit *Rabbit) PublishJSON(key string, body interface{}) error {
	dat, _ := json.Marshal(body)
	Error := rabbit.Chan.Publish(
		rabbit.ExName, // exchange
		key,           // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        dat,
		})

	return Error
}

func (rabbit *Rabbit) AnonymousQueue() err {
	var Error error

	rabbit.Queue, Error = rabbit.Chan(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return Error
}

func (rabbit *Rabbit) BindKey(... keys) error{
	
}
