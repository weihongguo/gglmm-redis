package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Cacher Redis缓存 实现了Cacher接口
type Cacher struct {
	pool      *redis.Pool
	expires   int
	keyPrefix string
	selfPool  bool
}

// NewCacherConfig --
func NewCacherConfig(config ConfigCacher) *Cacher {
	idelTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", config.IdelTimeout))
	if err != nil {
		log.Fatal(err)
	}
	return NewCacher(
		config.Network,
		config.Address,
		config.MaxActive,
		config.MaxIdel,
		idelTimeout,
		config.Expires,
	)
}

// NewCacher --
func NewCacher(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration, expires int) *Cacher {
	pool := NewPool(network, address, maxActive, maxIdle, idleTimeout)
	if expires < 10 {
		expires = 10
	}
	return &Cacher{
		pool:     pool,
		expires:  expires,
		selfPool: true,
	}
}

// NewCacherPool --
func NewCacherPool(pool *redis.Pool, expires int) *Cacher {
	if expires < 10 {
		expires = 10
	}
	return &Cacher{
		pool:     pool,
		expires:  expires,
		selfPool: false,
	}
}

// Close --
func (cacher *Cacher) Close() {
	if cacher.selfPool {
		cacher.pool.Close()
	}
}

// SetExpires --
func (cacher *Cacher) SetExpires(expires int) {
	if expires < 10 {
		expires = 10
	}
	cacher.expires = expires
}

// Expires --
func (cacher *Cacher) Expires() int {
	return cacher.expires
}

// SetKeyPrefix --
func (cacher *Cacher) SetKeyPrefix(keyPrefix string) {
	cacher.keyPrefix = keyPrefix
}

// KeyPrefix --
func (cacher *Cacher) KeyPrefix() string {
	return cacher.keyPrefix
}

// SetEx --
func (cacher *Cacher) SetEx(key string, value interface{}, ex int) error {
	conn := cacher.pool.Get()
	if conn == nil {
		return ErrConnect
	}
	defer conn.Close()

	key = cacher.KeyPrefix() + key

	var reply interface{}
	var err error
	switch reflect.TypeOf(value).Kind() {
	case reflect.Struct, reflect.Slice, reflect.Map, reflect.Ptr:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return err
		}
		reply, err = conn.Do("SET", key, jsonValue, "EX", ex)
	default:
		reply, err = conn.Do("SET", key, value, "EX", ex)
	}
	ok, err := redis.String(reply, err)
	if err != nil {
		return err
	}
	if ok != "OK" {
		return ErrReply
	}
	return nil
}

// Set --
func (cacher *Cacher) Set(key string, value interface{}) error {
	return cacher.SetEx(key, value, cacher.expires)
}

// Get --
func (cacher *Cacher) Get(key string) (interface{}, error) {
	conn := cacher.pool.Get()
	if conn == nil {
		return nil, ErrConnect
	}
	defer conn.Close()

	key = cacher.KeyPrefix() + key
	return conn.Do("GET", key)
}

// Del --
func (cacher *Cacher) Del(key string) (int, error) {
	conn := cacher.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	key = cacher.KeyPrefix() + key
	return redis.Int(conn.Do("DEL", key))
}

// DelPattern --
func (cacher *Cacher) DelPattern(pattern string) (int, error) {
	conn := cacher.pool.Get()
	if conn == nil {
		return 0, ErrConnect
	}
	defer conn.Close()

	pattern = cacher.KeyPrefix() + pattern
	keys, err := redis.Values(conn.Do("KEYS", pattern))
	if err != nil {
		return 0, err
	}
	return redis.Int(conn.Do("DEL", keys...))
}

// GetInt --
func (cacher *Cacher) GetInt(key string) (int, error) {
	return redis.Int(cacher.Get(key))
}

// GetInt64 --
func (cacher *Cacher) GetInt64(key string) (int64, error) {
	return redis.Int64(cacher.Get(key))
}

// GetFloat64 --
func (cacher *Cacher) GetFloat64(key string) (float64, error) {
	return redis.Float64(cacher.Get(key))
}

// GetBytes --
func (cacher *Cacher) GetBytes(key string) ([]byte, error) {
	return redis.Bytes(cacher.Get(key))
}

// GetString --
func (cacher *Cacher) GetString(key string) (string, error) {
	return redis.String(cacher.Get(key))
}

// GetObj --
func (cacher *Cacher) GetObj(key string, obj interface{}) error {
	value, err := cacher.GetBytes(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(value, obj)
}
