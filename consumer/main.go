package consumer

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/Zombispormedio/smartdb/lib/rabbit"
	"github.com/Zombispormedio/smartdb/lib/store"
	"gopkg.in/mgo.v2"
)

type Consumer struct {
	Client  *rabbit.Rabbit
	Keys    []string
	Routing func(*Deliver) error
	Session *mgo.Session
}

type Deliver struct {
	Consumer *Consumer
	Delivery rabbit.Delivery
}

func New(routing func(*Deliver) error, session *mgo.Session) (*Consumer, error) {
	var Error error
	consumer := Consumer{}
	consumer.Routing = routing
	consumer.Session = session

	rClient, RError := rabbit.New(os.Getenv("EX_RABBIT"), "topic", true)

	if RError != nil {
		return nil, RError
	}

	QueueError := rClient.AnonymousQueue()

	if QueueError != nil {
		return nil, QueueError
	}

	consumer.Client = rClient

	consumer.Bind()

	return &consumer, Error
}

func (consumer *Consumer) Bind() error {
	var Error error
	client := consumer.Client

	var firstKey string
	Error = store.Get("push_identifier", "Config", func(value string) {
		firstKey = value + ".#"
	})

	if Error != nil {
		return Error
	}

	consumer.AppendKey(firstKey)

	Error = client.BindKeys(consumer.Keys...)

	return Error
}

func (consumer *Consumer) ReBind() error {
	var Error error
	client := consumer.Client

	Error = client.UnbindKeys(consumer.Keys...)
	if Error != nil {
		return Error
	}
	consumer.Keys = nil

	Error = consumer.Bind()

	return Error
}

func (consumer *Consumer) AppendKey(key string) {
	consumer.Keys = append(consumer.Keys, key)
}

func (consumer *Consumer) Run() error {

	if consumer.Routing == nil {
		return errors.New("Error: Consumer don't have routing function. Consumer didn't start")
	}

	client := consumer.Client

	msgs, Error := client.Consume()
	if Error != nil {
		return Error
	}

	log.Info("Consumer started")

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			deliver := Deliver{}
			log.Info(string(msg.Body))
			deliver.Consumer = consumer
			deliver.Delivery = &msg
			RouteError := consumer.Routing(&deliver)

			if RouteError != nil {
				log.Error(RouteError)
			}

		}
	}()
	<-forever

	log.Error(Error)
	return Error
}

func (deliver *Deliver) Session() *mgo.Session {
	return deliver.Consumer.Session.Copy()
}

func (deliver *Deliver) Body(result *map[string]interface{}) error {

	var Error error

	delivery := deliver.Delivery

	Error = json.Unmarshal(delivery.Body, result)

	return Error
}

type KeyElem struct {
	Parent   string
	Children []string
}

func (deliver *Deliver) KeyElem() *KeyElem {
	delivery := deliver.Delivery
	keys := strings.Split(delivery.RoutingKey, ".")

	keyElem := KeyElem{}

	keyElem.Parent = keys[0]
	keyElem.Children = keys[1:]

	return &keyElem
}
