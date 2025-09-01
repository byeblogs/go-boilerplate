package dtos

type NSQHeader struct {
	RequestID string `json:"request_id,omitempty"`
	Topic     string `json:"topic,omitempty"`
	Service   string `json:"service,omitempty"`
}

type NSQRequest[T any] struct {
	Data    T         `json:"data,omitempty"`
	Headers NSQHeader `json:"headers,omitempty"`
}