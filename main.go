package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func die(err error) {
	if err != nil {
		log.Fatalf(err.Error())
		panic(fmt.Sprintf("%s", err))
	}
}

func main() {
	cfg, err := config()
	die(err)
	log.Printf("Config: %+v", cfg)

	conn, err := amqp.Dial(cfg.Rmq)
	die(err)
	defer conn.Close()

	ch, err := conn.Channel()
	die(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		cfg.Exchange, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	die(err)

	msgs1, err := registerInput(cfg.In1, ch, cfg)
	msgs2, err := registerInput(cfg.In2, ch, cfg)
	die(err)

	forever := make(chan bool)

	go func() {
		for {
			select {
			case in := <-msgs1:
				log.Printf("In1: %s", in.Body)

			case in := <-msgs2:
				log.Printf("In2: %s", in.Body)
			}
		}
	}()

	log.Printf("amqp.Delivery [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func registerInput(rk string, ch *amqp.Channel, cfg Cfg) (<-chan amqp.Delivery, error) {
	q, err := ch.QueueDeclare(
		rk,    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare a queue: %q", err)
	}

	log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, cfg.Exchange, q.Name)
	err = ch.QueueBind(
		q.Name,       // queue name
		q.Name,       // routing key
		cfg.Exchange, // exchange
		false,
		nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to bind a queue: %q", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to register a consumer: %q", err)
	}

	return msgs, nil
}
