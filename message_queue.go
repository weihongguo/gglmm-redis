package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// MessageQueue --
type MessageQueue struct {
	pool     *redis.Pool
	channel  string
	selfPool bool
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
	if channel == "" {
		return nil
	}
	pool := NewPool(network, address, maxActive, maxIdle, idleTimeout)
	return &MessageQueue{
		pool:     pool,
		channel:  channel,
		selfPool: true,
	}
}

// NewMessageQueuePool --
func NewMessageQueuePool(pool *redis.Pool, channel string) *MessageQueue {
	if channel == "" {
		return nil
	}
	return &MessageQueue{
		pool:     pool,
		channel:  channel,
		selfPool: false,
	}
}

// Close --
func (mq *MessageQueue) Close() {
	if mq.selfPool {
		mq.pool.Close()
	}
}

// Publish --
func (mq *MessageQueue) Publish(message []byte) error {
	if message == nil || len(message) == 0 {
		return ErrChannelEmpty
	}

	conn := mq.pool.Get()
	if conn == nil {
		return ErrConnect
	}
	defer conn.Close()

	_, err := conn.Do("LPUSH", mq.channel, message)
	return err
}

// Subscribe --
func (mq *MessageQueue) Subscribe(handler func(message []byte, err error), timeout int) (chan<- interface{}, error) {
	conn := mq.pool.Get()
	if conn == nil {
		return nil, ErrConnect
	}

	if timeout < 10 {
		timeout = 10
	}

	stop := false
	go func() {
		defer conn.Close()
		for {
			if stop {
				log.Printf("[%s] message queue stop!\n", mq.channel)
				return
			}
			reply, err := conn.Do("BRPOP", mq.channel, timeout)
			if err != nil {
				handler(nil, err)
			} else {
				messages, err := redis.ByteSlices(reply, err)
				if err != nil {
					handler(nil, err)
				} else {
					if len(messages) != 2 && string(messages[0]) != mq.channel {
						handler(nil, ErrChannel)
					} else {
						handler(messages[1], nil)
					}
				}
			}
		}
	}()

	stopChan := make(chan interface{})
	go func() {
		defer close(stopChan)
		<-stopChan
		log.Printf("[%s] message queue will stop!\n", mq.channel)
		stop = true
	}()

	return stopChan, nil
}
