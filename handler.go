package main

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type product struct {
	m1 json.RawMessage
	m2 json.RawMessage
}

type pid int

type products map[pid]product

func handleInput1(m amqp.Delivery) {}

func handleInput2(m amqp.Delivery) {}
