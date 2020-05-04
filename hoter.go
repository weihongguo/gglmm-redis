package redis

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Hoter --
type Hoter struct {
	pool     *redis.Pool
	name     string
	selfPool bool
}

// NewHoterConfig --
func NewHoterConfig(config ConfigCounter, name string) *Hoter {
	if !config.Check() {
		log.Printf("%+v\n", config)
		log.Fatal("ConfigCounter invalid")
	}
	idelTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", config.IdelTimeout))
	if err != nil {
		log.Fatal(err)
	}
	return NewHoter(
		config.Network,
		config.Address,
		config.MaxActive,
		config.MaxIdel,
		idelTimeout,
		name,
	)
}

// NewHoter --
func NewHoter(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration, name string) *Hoter {
	pool := NewPool(network, address, maxActive, maxIdle, idleTimeout)
	if name == "" {
		return nil
	}
	return &Hoter{
		pool:     pool,
		name:     name,
		selfPool: true,
	}
}

// NewHoterPool --
func NewHoterPool(pool *redis.Pool, name string) *Hoter {
	if name == "" {
		return nil
	}
	return &Hoter{
		pool:     pool,
		name:     name,
		selfPool: false,
	}
}

// Close --
func (hoter *Hoter) Close() {
	if hoter.selfPool {
		hoter.pool.Close()
	}
}

// Del --
func (hoter *Hoter) Del() (int, error) {
	conn := hoter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("DEL", hoter.name))
}

// Add --
func (hoter *Hoter) Add(score, id int64) (int, error) {
	conn := hoter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("ZADD", hoter.name, score, id))
}

// Card --
func (hoter *Hoter) Card() (int, error) {
	conn := hoter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("ZCARD", hoter.name))
}

// Count --
func (hoter *Hoter) Count(min, max int64) (int, error) {
	conn := hoter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("ZCOUNT", hoter.name, min, max))
}

// Rem --
func (hoter *Hoter) Rem(id int64) (int, error) {
	conn := hoter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("ZREM", hoter.name, id))
}

// RemRangeByScore --
func (hoter *Hoter) RemRangeByScore(min, max int64) (int, error) {
	conn := hoter.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	return redis.Int(conn.Do("ZREMRANGEBYSCORE", hoter.name, min, max))
}

// RangeByScore --
func (hoter *Hoter) RangeByScore(max int64) ([]int64, []int64, error) {
	conn := hoter.pool.Get()
	if conn == nil {
		return nil, nil, ErrConnect
	}
	defer conn.Close()

	results, err := redis.Int64s(conn.Do("ZRANGEBYSCORE", hoter.name, math.MinInt64, max, "WITHSCORES"))
	if err != nil {
		return nil, nil, err
	}
	scores := make([]int64, 0)
	ids := make([]int64, 0)
	for i, result := range results {
		if (i % 2) == 0 {
			ids = append(ids, result)
		} else {
			scores = append(scores, result)
		}
	}
	return ids, scores, nil
}

// RevRangeByScore --
func (hoter *Hoter) RevRangeByScore(min int64) ([]int64, []int64, error) {
	conn := hoter.pool.Get()
	if conn == nil {
		return nil, nil, ErrConnect
	}
	defer conn.Close()

	results, err := redis.Int64s(conn.Do("ZREVRANGEBYSCORE", hoter.name, math.MaxInt64, min, "WITHSCORES"))
	if err != nil {
		return nil, nil, err
	}
	scores := make([]int64, 0)
	ids := make([]int64, 0)
	for i, result := range results {
		if (i % 2) == 0 {
			ids = append(ids, result)
		} else {
			scores = append(scores, result)
		}
	}
	return ids, scores, nil
}
