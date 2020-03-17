package redis

import (
	"errors"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// NewRedisPool --
func NewRedisPool(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration) *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(network, address)
		},
	}

	if err := Ping(redisPool); err != nil {
		log.Fatal(err)
	}

	return redisPool
}

// Ping --
func Ping(pool *redis.Pool) error {
	redisConn := pool.Get()
	if redisConn == nil {
		return ErrConn
	}
	defer redisConn.Close()

	pong, err := redisConn.Do("PING")
	if err != nil {
		return err
	}
	pong, err = redis.String(pong, err)
	if err != nil {
		return err
	}
	if pong != "PONG" {
		return errors.New("PING NOT PONG")
	}
	return nil
}
