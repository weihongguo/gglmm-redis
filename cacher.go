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
	redisPool *redis.Pool
	expires   int
	keyPrefix string
	name      string
}

// NewCacherConfig --
func NewCacherConfig(config ConfigCacher) *Cacher {
	if !config.Check() {
		log.Printf("%+v\n", config)
		log.Fatal("ConfigCacher invalid")
	}
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
	redisPool := NewRedisPool(network, address, maxActive, maxIdle, idleTimeout)

	if expires < 10 {
		expires = 10
	}

	return &Cacher{
		redisPool: redisPool,
		expires:   expires,
		name:      "redis",
	}
}

// Close --
func (cacher *Cacher) Close() {
	cacher.redisPool.Close()
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

// SetName --
func (cacher *Cacher) SetName(name string) {
	cacher.name = name
}

// Name --
func (cacher *Cacher) Name() string {
	return cacher.name
}

// SetEx --
func (cacher *Cacher) SetEx(key string, value interface{}, ex int) error {
	redisConn := cacher.redisPool.Get()
	if redisConn == nil {
		return ErrConn
	}
	defer redisConn.Close()

	key = cacher.KeyPrefix() + key

	var err error
	switch reflect.TypeOf(value).Kind() {
	case reflect.Struct, reflect.Slice, reflect.Map, reflect.Ptr:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return err
		}
		_, err = redisConn.Do("SET", key, jsonValue, "EX", ex)
	default:
		_, err = redisConn.Do("SET", key, value, "EX", ex)
	}
	return err
}

// Set --
func (cacher *Cacher) Set(key string, value interface{}) error {
	return cacher.SetEx(key, value, cacher.expires)
}

// Get --
func (cacher *Cacher) Get(key string) (interface{}, error) {
	redisConn := cacher.redisPool.Get()
	if redisConn == nil {
		return nil, ErrConn
	}
	defer redisConn.Close()

	key = cacher.KeyPrefix() + key
	return redisConn.Do("GET", key)
}

// Del --
func (cacher *Cacher) Del(key string) error {
	redisConn := cacher.redisPool.Get()
	if redisConn == nil {
		return ErrConn
	}
	defer redisConn.Close()

	_, err := redisConn.Do("DEL", key)
	return err
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
