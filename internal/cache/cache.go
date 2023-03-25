package cache

import (
	"fmt"
	"time"

	"github.com/programzheng/black-key/config"
)

type CacheDriver interface {
	HGet(key string, field string) (string, error)
	HDel(key string, fields ...string) (int64, error)
	Del(keys ...string) (int64, error)
	HSet(key string, values ...interface{}) (int64, error)
	Exists(keys ...string) (int64, error)
	SIsMember(key string, member interface{}) (bool, error)
	SMembers(key string) ([]string, error)
	SRem(key string, members ...interface{}) (int64, error)
	SAdd(key string, members ...interface{}) (int64, error)
	Expire(key string, expiration time.Duration) (bool, error)
}

func GetCacheDriver(driverName string) (CacheDriver, error) {
	if driverName == "" {
		driverName = config.Cfg.GetString("CACHE_DRIVER")
	}
	switch driverName {
	case "redis":
		return newRedisCacheDriver("", "", false, false, 0), nil
	default:
		return nil, fmt.Errorf("unsupported cache driver: %s", driverName)
	}
}
