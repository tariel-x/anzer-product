package main

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type product struct {
	m1 json.RawMessage
	m2 json.RawMessage
}

type pid string

type products map[pid]product

const (
	pidHeader = "pid"
)

var (
	storage = products{}
)

func handleInput1(m amqp.Delivery) {
	tmpID, ok := m.Headers[pidHeader]
	if !ok {
		return
	}
	id := tmpID.(pid)
	if _, ok := storage[id]; !ok {
		storage[id] = product{}
	}

	product := storage[id]
	product.m1 = json.RawMessage(m.Body)

	data := makeProduct(id)
	if data != nil {
		push(data, id, "")
		clear(id)
	}
}

func handleInput2(m amqp.Delivery) {
	id, ok := m.Headers[pidHeader]
	if !ok {
		return
	}
	if _, ok := storage[id.(pid)]; !ok {
		storage[id.(pid)] = product{}
	}

	product := storage[id.(pid)]
	product.m2 = json.RawMessage(m.Body)

	data := makeProduct(id.(pid))
	if data != nil {
		push(data)
		clear(id.(pid))
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

	productData := map[string]json.RawMessage{
		"m1": product.m1,
		"m2": product.m2,
	}

	data, _ := json.Marshal(productData)
	return data
}

func clear(id pid) {
	if _, ok := storage[id]; !ok {
		return
	}
	delete(storage, id)
}

func push(data json.RawMessage, id pid, rk string) {}
