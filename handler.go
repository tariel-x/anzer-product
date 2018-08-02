package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type product struct {
	m1 interface{}
	m2 interface{}
}

type pid string

type products map[pid]product

const (
	pidHeader = "pid"
)

var (
	storage = products{}
)

func handleInput1(m amqp.Delivery, ch *amqp.Channel) {
	tmpID, ok := m.Headers[pidHeader]
	if !ok {
		return
	}
	id := pid(tmpID.(string))
	if _, ok := storage[id]; !ok {
		storage[id] = product{}
	}

	log.Printf("Add first part to pid %q", id)
	product := storage[id]
	err := json.Unmarshal(m.Body, &product.m1)
	if err != nil {
		log.Printf("Error unmarshalling data %q: %q", m.Body, err)
		return
	}
	storage[id] = product
	log.Printf("Current product is %+v", storage[id])

	data := makeProduct(id)
	if data != nil {
		push(data, id, ch)
		clear(id)
	}
}

func handleInput2(m amqp.Delivery, ch *amqp.Channel) {
	tmpID, ok := m.Headers[pidHeader]
	if !ok {
		return
	}
	id := pid(tmpID.(string))
	if _, ok := storage[id]; !ok {
		storage[id] = product{}
	}

	log.Printf("Add second part to pid %q", id)
	product := storage[id]
	err := json.Unmarshal(m.Body, &product.m2)
	if err != nil {
		log.Printf("Error unmarshalling data %q: %q", m.Body, err)
		return
	}
	storage[id] = product
	log.Printf("Current product is %+v", storage[id])

	data := makeProduct(id)
	if data != nil {
		push(data, id, ch)
		clear(id)
	}
}

func makeProduct(id pid) json.RawMessage {
	product, ok := storage[id]
	if !ok {
		return nil
	}
	if product.m1 == nil || product.m2 == nil {
		return nil
	}

	productData := map[string]interface{}{
		Type1: product.m1,
		Type2: product.m2,
	}

	data, err := json.Marshal(productData)
	if err != nil {
		log.Printf("Error marshalling product data: %q", err)
	}
	return data
}

func clear(id pid) {
	if _, ok := storage[id]; !ok {
		return
	}
	delete(storage, id)
}

func push(data json.RawMessage, id pid, ch *amqp.Channel) {
	log.Printf("Sending product pid %q message further to %q", id, Out)
	err := ch.Publish(
		Exchange, // exchange
		Out,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Headers: amqp.Table{
				pidHeader: string(id),
			},
			Body: []byte(data),
		})
	if err != nil {
		log.Printf("Error while sending message further: %q", err)
	}
}
