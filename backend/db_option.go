package backend

import "time"

const (
	defaultConnTimeout = 2 * time.Second
	defaultIdleTimeout = 10 * 60 * time.Second

	DefaultCapacity    = 128
	DefaultMaxCapacity = 1024
)

var (
	DefaultOption = &Option{
		idleTimeout: defaultIdleTimeout,
		connTimeout: defaultConnTimeout,
		maxCapacity: DefaultMaxCapacity,
		capacity:    DefaultCapacity,
	}
)

type Option struct {
	maxCapacity int
	capacity    int
	connTimeout time.Duration
	idleTimeout time.Duration
}

type DbOption func(o *Option)

func WithCapacity(capacity int) DbOption {
	return func(o *Option) {
		o.capacity = capacity
	}
}

func WithMaxCapacity(maxCapacity int) DbOption {
	return func(o *Option) {
		o.maxCapacity = maxCapacity
	}
}

func WithConnTimeout(connTimeout time.Duration) DbOption {
	return func(o *Option) {
		o.connTimeout = connTimeout
	}
}

func WithIdleTimeout(idleTimeout time.Duration) DbOption {
	return func(o *Option) {
		o.idleTimeout = idleTimeout
	}
}
