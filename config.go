package main

import (
	"errors"
	"os"
)

type Cfg struct {
	Rmq      string
	In1      string
	In2      string
	Out      string
	Exchange string
}

func config() (Cfg, error) {
	cfg := Cfg{
		Exchange: "anzer",
	}
	cfg.Rmq = os.Getenv("RMQ")
	cfg.In1 = os.Getenv("IN1")
	cfg.In2 = os.Getenv("IN2")
	cfg.Out = os.Getenv("OUT")

	if cfg.Rmq == "" || cfg.In1 == "" || cfg.In2 == "" || cfg.Out == "" {
		return cfg, errors.New("Some envs are missing")
	}

	return cfg, nil
}
