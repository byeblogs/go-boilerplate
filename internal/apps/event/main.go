package main

import (
	"github.com/joho/godotenv"
	"github.com/nsqio/go-nsq"
)

func main() {
	_ = godotenv.Load(".env")
	
	consumers := []*nsq.Consumer{}

	module := Module{}
	module.New(consumers)

	for _, consumer := range consumers {
		go consumer.Stop()
	}
}