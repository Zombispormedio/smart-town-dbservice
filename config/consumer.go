package config

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"github.com/Zombispormedio/smartdb/lib/rabbit"
	"github.com/Zombispormedio/smartdb/lib/store"
)

type Consumer struct {
	Client *rabbit.Rabbit
    Keys   []string
    
}

func CreateConsumer() (*Consumer, error) {
	var Error error
	consumer := Consumer{}

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
		firstKey = value+".#"
	})
    
    if Error != nil{
        return Error
    }
    
    consumer.AppendKey(firstKey)
    
    Error=client.BindKeys(consumer.Keys...)
    
	return Error
}

func (consumer *Consumer) ReBind() error {
    var Error error
	client := consumer.Client
    
    Error=client.UnbindKeys(consumer.Keys...)
     if Error != nil{
        return Error
    }
    consumer.Keys=nil
    
    Error=consumer.Bind()
      
    return Error
}


func (consumer *Consumer) AppendKey(key string) {
    consumer.Keys=append(consumer.Keys, key)
}


func (consumer *Consumer) Run() error{
    log.Info("go run")
    client := consumer.Client
    
    msgs, Error:=client.Consume()
    if Error != nil{
        return Error
    }
    
    forever := make(chan bool)
    go func(){
        for msg :=range msgs{
            log.Info(msg)
        }
    }()
    <-forever
    
    log.Error(Error)
    return Error
}