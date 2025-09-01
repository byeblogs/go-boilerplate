package subscribers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-boilerplate/internal/constants"
	"go-boilerplate/internal/dtos"
	"go-boilerplate/internal/usecases"

	"github.com/nsqio/go-nsq"
	"github.com/sarulabs/di"
)

type UserSubscriber interface {
	GetUserByID(message *nsq.Message) error
}

type UserSubscriberImpl struct {
	usecase *usecases.Usecase
}

func NewUserSubscriber(ioc di.Container) UserSubscriber {
	return &UserSubscriberImpl{
		usecase: ioc.Get(constants.Usecase).(*usecases.Usecase),
	}
}

func (s *UserSubscriberImpl) GetUserByID(message *nsq.Message) error {
	var (
		params   = dtos.NSQRequest[dtos.UserIDParams]{}
		ctx      = context.TODO()
		maxRetry = uint16(len(constants.NSQRetryDelayMap))
		err      error
	)

	if err = json.Unmarshal(message.Body, &params); err != nil {
		return err
	}

	fmt.Printf("params: %+v\n", params)

	ctx = context.WithValue(ctx, constants.TracerRequestID, params.Headers.RequestID)
	// err = h.usecase.Voucher.RefundCoin(ctx, &params.Data)
	if err != nil {
		if message.Attempts <= maxRetry {
			message.RequeueWithoutBackoff(constants.NSQRetryDelayMap[message.Attempts])
		} else {
			return err
		}
	}

	message.Finish()
	return nil
}
