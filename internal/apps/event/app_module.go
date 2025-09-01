package main

import (
	"go-boilerplate/internal/apps"

	"github.com/nsqio/go-nsq"
)

type Module struct{}

func (m *Module) New(consumers []*nsq.Consumer) {
	ioc := apps.NewIOC()

	registry := NewRegistry(consumers, ioc)
	registry.Init()
}