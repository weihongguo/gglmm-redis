package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// NewPool --
func NewPool(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(network, address)
		},
	}
	if err := Ping(pool); err != nil {
		log.Fatal(err)
	}
	return pool
}

// Ping --
func Ping(pool *redis.Pool) error {
	conn := pool.Get()
	if conn == nil {
		return ErrConnect
	}
	defer conn.Close()

	pong, err := conn.Do("PING")
	if err != nil {
		return err
	}
	pong, err = redis.String(pong, err)
	if err != nil {
		return err
	}
	if pong != "PONG" {
		return ErrPing
	}
	return nil
}

// Factory --
type Factory struct {
	pool *redis.Pool
}

// NewFactoryConfig --
func NewFactoryConfig(config ConfigRedis) *Factory {
	if !config.Check() {
		log.Printf("%+v\n", config)
		log.Fatal("ConfigRedis invalid")
	}
	idelTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", config.IdelTimeout))
	if err != nil {
		log.Fatal(err)
	}
	return NewFactory(
		config.Network,
		config.Address,
		config.MaxActive,
		config.MaxIdel,
		idelTimeout,
	)
}

// NewFactory --
func NewFactory(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration) *Factory {
	pool := NewPool(network, address, maxActive, maxIdle, idleTimeout)
	return &Factory{
		pool: pool,
	}
}

// Close --
func (factory *Factory) Close() {
	factory.pool.Close()
}

// NewCacher --
func (factory *Factory) NewCacher(expires int) *Cacher {
	return NewCacherPool(
		factory.pool,
		expires,
	)
}

// NewCounter --
func (factory *Factory) NewCounter(name string) *Counter {
	return NewCounterPool(
		factory.pool,
		name,
	)
}

// NewToper --
func (factory *Factory) NewToper(name string, limit int) *Toper {
	return NewToperPool(
		factory.pool,
		name,
		limit,
	)
}

// NewHoter --
func (factory *Factory) NewHoter(name string) *Hoter {
	return NewHoterPool(
		factory.pool,
		name,
	)
}

// NewMessageQueue --
func (factory *Factory) NewMessageQueue(channel string) *MessageQueue {
	return NewMessageQueuePool(
		factory.pool,
		channel,
	)
}
