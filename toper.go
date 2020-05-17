package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Toper --
type Toper struct {
	pool     *redis.Pool
	name     string
	limit    int
	selfPool bool
}

// NewToperConfig --
func NewToperConfig(config ConfigToper, name string, limit int) *Toper {
	if !config.Check() {
		log.Printf("%+v\n", config)
		log.Fatal("ConfigCounter invalid")
	}
	idelTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", config.IdelTimeout))
	if err != nil {
		log.Fatal(err)
	}
	return NewToper(
		config.Network,
		config.Address,
		config.MaxActive,
		config.MaxIdel,
		idelTimeout,
		name,
		limit,
	)
}

// NewToper --
func NewToper(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration, name string, limit int) *Toper {
	pool := NewPool(network, address, maxActive, maxIdle, idleTimeout)
	if name == "" {
		return nil
	}
	return &Toper{
		pool:     pool,
		name:     name,
		limit:    limit,
		selfPool: true,
	}
}

// NewToperPool --
func NewToperPool(pool *redis.Pool, name string, limit int) *Toper {
	if name == "" {
		return nil
	}
	return &Toper{
		pool:     pool,
		name:     name,
		limit:    limit,
		selfPool: false,
	}
}

// Close --
func (toper *Toper) Close() {
	if toper.selfPool {
		toper.pool.Close()
	}
}

// Del --
func (toper *Toper) Del() error {
	conn := toper.pool.Get()
	if conn == nil {
		return ErrConnect
	}
	defer conn.Close()

	_, err := conn.Do("DEL", toper.name)
	return err
}

// Push --
func (toper *Toper) Push(id int64) error {
	conn := toper.pool.Get()
	if conn == nil {
		return ErrConnect
	}
	defer conn.Close()

	_, err := conn.Do("LPUSH", toper.name, id)
	if err != nil {
		return err
	}
	_, err = conn.Do("LTRIM", toper.name, 0, toper.limit-1)
	return err
}

// Pop --
func (toper *Toper) Pop() (int64, error) {
	conn := toper.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int64(conn.Do("LPOP", toper.name))
}

// Range --
func (toper *Toper) Range() ([]int64, error) {
	conn := toper.pool.Get()
	if conn == nil {
		return nil, ErrConnect
	}
	defer conn.Close()

	return redis.Int64s(conn.Do("LRANGE", toper.name, 0, -1))
}
