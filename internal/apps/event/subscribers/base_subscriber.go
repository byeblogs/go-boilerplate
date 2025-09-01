package subscribers

import (
	"github.com/sarulabs/di"
)

type Subscriber struct {
	User UserSubscriber
}

func NewSubscriber(ioc di.Container) *Subscriber {
	return &Subscriber{
		User: NewUserSubscriber(ioc),
	}
}