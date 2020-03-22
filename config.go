package redis

import (
	"errors"
	"log"
	"strings"
)

var (
	// ErrConn --
	ErrConn = errors.New("连接错误")
)

// ConfigRedis --
type ConfigRedis struct {
	Network     string
	Address     string
	MaxActive   int
	MaxIdel     int
	IdelTimeout int
}

// Check --
func (config ConfigRedis) Check() bool {
	if config.Network == "" {
		return false
	}
	if config.Address == "" || !strings.Contains(config.Address, ":") {
		return false
	}
	if config.MaxActive <= 0 || config.MaxIdel <= 0 || config.IdelTimeout <= 0 {
		return false
	}
	log.Println("config redis check valid")
	return true
}

// ConfigCacher --
type ConfigCacher struct {
	ConfigRedis
	Expires int
}

// ConfigCounter --
type ConfigCounter struct {
	ConfigRedis
}

// ConfigToper --
type ConfigToper struct {
	ConfigRedis
}

// ConfigHoter --
type ConfigHoter struct {
	ConfigRedis
}

// ConfigMessageQueue --
type ConfigMessageQueue struct {
	ConfigRedis
}
