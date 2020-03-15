package redis

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// MessageQueue --
type MessageQueue struct {
	redisPool *redis.Pool
	channel   string
	name      string
}

// NewMessageQueueConfig --
func NewMessageQueueConfig(config ConfigMessageQueue, channel string) *MessageQueue {
	if !config.Check() {
		log.Printf("%+v\n", config)
		log.Fatal("ConfigMessageQueue invalid")
	}
	idelTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", config.IdelTimeout))
	if err != nil {
		log.Fatal(err)
	}
	return NewMessageQueue(
		config.Network,
		config.Address,
		config.MaxActive,
		config.MaxIdel,
		idelTimeout,
		channel,
	)
}

// NewMessageQueue --
func NewMessageQueue(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration, channel string) *MessageQueue {
	return &MessageQueue{
		redisPool: &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTimeout * time.Second,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				return redis.Dial(network, address)
			},
		},
		channel: channel,
		name:    "redis",
	}
}

// Close --
func (mq *MessageQueue) Close() {
	mq.redisPool.Close()
}

// SetName --
func (mq *MessageQueue) SetName(name string) {
	mq.name = name
}

// Name --
func (mq *MessageQueue) Name() string {
	return mq.name
}

// Push --
func (mq *MessageQueue) Push(message []byte) error {
	if message == nil || len(message) == 0 {
		return errors.New("channel is empty")
	}

	redisConn := mq.redisPool.Get()
	if redisConn == nil {
		return ErrConn
	}
	defer redisConn.Close()

	_, err := redisConn.Do("lpush", mq.channel, message)
	return err
}

// BPop --
func (mq *MessageQueue) BPop(handler func(message []byte, err error), timeout int) (chan<- interface{}, error) {
	redisConn := mq.redisPool.Get()
	if redisConn == nil {
		return nil, ErrConn
	}

	if timeout < 10 {
		timeout = 10
	}

	stop := false
	go func() {
		defer redisConn.Close()
		for {
			if stop {
				log.Printf("[%s] message queue stop!\n", mq.channel)
				return
			}
			reply, err := redisConn.Do("brpop", mq.channel, timeout)
			if err != nil {
				handler(nil, err)
			} else {
				messages, err := redis.ByteSlices(reply, err)
				if err != nil {
					handler(nil, err)
				} else {
					if len(messages) != 2 && string(messages[0]) != mq.channel {
						handler(nil, errors.New("channel error"))
					} else {
						handler(messages[1], nil)
					}
				}
			}
		}
	}()

	stopChan := make(chan interface{})
	go func() {
		<-stopChan
		log.Printf("[%s] message queue will stop!\n", mq.channel)
		stop = true
	}()

	return stopChan, nil
}
