package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Counter --
type Counter struct {
	redisPool *redis.Pool
	name      string
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
	redisPool := NewRedisPool(network, address, maxActive, maxIdle, idleTimeout)
	return &Counter{
		redisPool: redisPool,
		name:      name,
	}
}

// Get --
func (counter *Counter) Get() (int, error) {
	redisConn := counter.redisPool.Get()
	if redisConn == nil {
		return 0, ErrConn
	}
	defer redisConn.Close()

	return redis.Int(redisConn.Do("GET", counter.name))
}

// Set --
func (counter *Counter) Set(value int) (int, error) {
	redisConn := counter.redisPool.Get()
	if redisConn == nil {
		return 0, ErrConn
	}
	defer redisConn.Close()

	oldValue, err := redis.Int(redisConn.Do("GET", counter.name))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			oldValue = 0
		} else {
			return 0, err
		}
	}

	_, err = redisConn.Do("SET", counter.name, value)
	return oldValue, err
}

// Zero --
func (counter *Counter) Zero() (int, error) {
	return counter.Set(0)
}

// Increase --
func (counter *Counter) Increase() (int, error) {
	redisConn := counter.redisPool.Get()
	if redisConn == nil {
		return 0, ErrConn
	}
	defer redisConn.Close()

	return redis.Int(redisConn.Do("INCR", counter.name))
}

// Decrease --
func (counter *Counter) Decrease() (int, error) {
	redisConn := counter.redisPool.Get()
	if redisConn == nil {
		return 0, ErrConn
	}
	defer redisConn.Close()

	return redis.Int(redisConn.Do("DECR", counter.name))
}

// IncreaseBy --
func (counter *Counter) IncreaseBy(diff int) (int, error) {
	redisConn := counter.redisPool.Get()
	if redisConn == nil {
		return 0, ErrConn
	}
	defer redisConn.Close()

	return redis.Int(redisConn.Do("INCRBY", counter.name, diff))
}

// DecreaseBy --
func (counter *Counter) DecreaseBy(diff int) (int, error) {
	redisConn := counter.redisPool.Get()
	if redisConn == nil {
		return 0, ErrConn
	}
	defer redisConn.Close()

	return redis.Int(redisConn.Do("DECRBY", counter.name, diff))
}
