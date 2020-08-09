package redis

import (
	"errors"
	"log"
	"strings"
)

var (
	// ErrConnect --
	ErrConnect = errors.New("连接错误")
	// ErrPing --
	ErrPing = errors.New("连接测试错误")
	// ErrReply --
	ErrReply = errors.New("回复错误")
	// ErrChannel --
	ErrChannel = errors.New("通道错误")
	// ErrChannelEmpty --
	ErrChannelEmpty = errors.New("通道为空")
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
	log.Println("ConfigRedis check pass")
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
