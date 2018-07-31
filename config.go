package main

import (
	"errors"
	"os"
)

func config() error {
	Exchange = "anzer"

	Rmq = os.Getenv("RMQ")
	In1 = os.Getenv("IN1")
	In2 = os.Getenv("IN2")
	Out = os.Getenv("OUT")
	Type1 = os.Getenv("TYPE1")
	Type2 = os.Getenv("TYPE2")

	if Rmq == "" || In1 == "" || In2 == "" || Out == "" || Type1 == "" || Type2 == "" {
		return errors.New("Some envs are missing")
	}

	return nil
}
