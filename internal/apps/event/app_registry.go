package main

import (
	"go-boilerplate/internal/apps/event/subscribers"
	"go-boilerplate/internal/constants"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/sarulabs/di"
)

type Registry struct {
	consumers  []*nsq.Consumer
	subscriber *subscribers.Subscriber
}

func NewRegistry(consumers []*nsq.Consumer, ioc di.Container) *Registry {
	return &Registry{
		consumers: consumers,
		subscriber: ioc.Get(constants.Subscriber).(*subscribers.Subscriber),
	}
}

func (r *Registry) Init() {
	config := nsq.NewConfig()
	config.MaxInFlight = 50
	config.LookupdPollInterval = 1 * time.Minute
}