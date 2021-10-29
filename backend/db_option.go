package backend

import (
	"github.com/fengleng/go-mysql-client/mysql"
	"time"
)

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
		charset:     mysql.DEFAULT_CHARSET,
		collationId: mysql.DEFAULT_COLLATION_ID,
	}
)

type Option struct {
	maxCapacity int
	capacity    int
	connTimeout time.Duration
	idleTimeout time.Duration

	collationId mysql.CollationId
	charset     string
}

type DbOption func(o *Option)

func WithCapacity(capacity int) DbOption {
	return func(o *Option) {
		o.capacity = capacity
	}
}

func WithCollationId(ci mysql.CollationId) DbOption {
	return func(o *Option) {
		o.collationId = ci
	}
}

func WithCharSet(charset string) DbOption {
	return func(o *Option) {
		o.charset = charset
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
