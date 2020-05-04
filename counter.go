package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Counter --
type Counter struct {
	pool     *redis.Pool
	name     string
	selfPool bool
}

// NewCounterConfig --
func NewCounterConfig(config ConfigCounter, name string) *Counter {
	if !config.Check() {
		log.Printf("%+v\n", config)
		log.Fatal("ConfigCounter invalid")
	}
	idelTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", config.IdelTimeout))
	if err != nil {
		log.Fatal(err)
	}
	return NewCounter(
		config.Network,
		config.Address,
		config.MaxActive,
		config.MaxIdel,
		idelTimeout,
		name,
	)
}

// NewCounter --
func NewCounter(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration, name string) *Counter {
	pool := NewPool(network, address, maxActive, maxIdle, idleTimeout)
	if name == "" {
		return nil
	}
	return &Counter{
		pool:     pool,
		name:     name,
		selfPool: true,
	}
}

// NewCounterPool --
func NewCounterPool(pool *redis.Pool, name string) *Counter {
	if name == "" {
		return nil
	}
	return &Counter{
		pool:     pool,
		name:     name,
		selfPool: false,
	}
}

// Close --
func (counter *Counter) Close() {
	if counter.selfPool {
		counter.pool.Close()
	}
}

// Del --
func (counter *Counter) Del() (int, error) {
	conn := counter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("DEL", counter.name))
}

// Get --
func (counter *Counter) Get() (int, error) {
	conn := counter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("GET", counter.name))
}

// Set --
func (counter *Counter) Set(value int) (int, error) {
	conn := counter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	oldValue, err := redis.Int(conn.Do("GET", counter.name))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			oldValue = 0
		} else {
			return 0, err
		}
	}

	reply, err := redis.String(conn.Do("SET", counter.name, value))
	if err != nil {
		return 0, err
	}
	if reply != "OK" {
		return 0, ErrReply
	}

	return oldValue, nil
}

// Zero --
func (counter *Counter) Zero() (int, error) {
	return counter.Set(0)
}

// Incr --
func (counter *Counter) Incr() (int, error) {
	conn := counter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("INCR", counter.name))
}

// Decr --
func (counter *Counter) Decr() (int, error) {
	conn := counter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("DECR", counter.name))
}

// IncrBy --
func (counter *Counter) IncrBy(diff int) (int, error) {
	conn := counter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("INCRBY", counter.name, diff))
}

// DecrBy --
func (counter *Counter) DecrBy(diff int) (int, error) {
	conn := counter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("DECRBY", counter.name, diff))
}
