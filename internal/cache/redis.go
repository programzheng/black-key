package cache

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/programzheng/black-key/config"

	"github.com/go-redis/redis/v8"
)

type RedisCacheDriver struct {
	client *redis.Client
}

func newRedisCacheDriver(addr string, password string, tls bool, tlsSkipVerify bool, db int) *RedisCacheDriver {
	if addr == "" {
		addr = config.Cfg.GetString("REDIS_ADDR")
	}
	if password == "" {
		password = config.Cfg.GetString("REDIS_PASSWORD")
	}
	if !tls {
		tls = config.Cfg.GetBool("REDIS_TLS")
	}
	if !tlsSkipVerify {
		tlsSkipVerify = config.Cfg.GetBool("REDIS_TLS_SKIP_VERIFY")
	}
	if db == 0 {
		db = config.Cfg.GetInt("REDIS_DB")
	}
	return &RedisCacheDriver{
		client: redis.NewClient(&redis.Options{
			Addr:      addr,
			TLSConfig: getTLSConfig(tls, tlsSkipVerify),
			Password:  password,
			DB:        db,
		}),
	}
}

func getTLSConfig(tlsBool bool, tlsSkipVerifyBool bool) *tls.Config {
	if tlsBool {
		tlsConfig := &tls.Config{}
		if tlsSkipVerifyBool {
			tlsConfig.InsecureSkipVerify = true
			return tlsConfig
		}
		tlsConfig.MinVersion = tls.VersionTLS12
		return tlsConfig
	}
	return nil
}

func (rcd *RedisCacheDriver) HGet(key string, field string) (string, error) {
	result, err := rcd.client.HGet(context.Background(), key, field).Result()
	return result, err
}

func (rcd *RedisCacheDriver) HDel(key string, fields ...string) (int64, error) {
	result, err := rcd.client.HDel(context.Background(), key, fields...).Result()
	return result, err
}

func (rcd *RedisCacheDriver) Del(keys ...string) (int64, error) {
	result, err := rcd.client.Del(context.Background(), keys...).Result()
	return result, err
}

func (rcd *RedisCacheDriver) HSet(key string, values ...interface{}) (int64, error) {
	result, err := rcd.client.HSet(context.Background(), key, values...).Result()
	return result, err
}

func (rcd *RedisCacheDriver) Exists(keys ...string) (int64, error) {
	result, err := rcd.client.Exists(context.Background(), keys...).Result()
	return result, err
}

func (rcd *RedisCacheDriver) SIsMember(key string, member interface{}) (bool, error) {
	result, err := rcd.client.SIsMember(context.Background(), key, member).Result()
	return result, err
}

func (rcd *RedisCacheDriver) SMembers(key string) ([]string, error) {
	result, err := rcd.client.SMembers(context.Background(), key).Result()
	return result, err
}

func (rcd *RedisCacheDriver) SRem(key string, members ...interface{}) (int64, error) {
	result, err := rcd.client.SRem(context.Background(), key, members...).Result()
	return result, err
}

func (rcd *RedisCacheDriver) SAdd(key string, members ...interface{}) (int64, error) {
	result, err := rcd.client.SAdd(context.Background(), key, members...).Result()
	return result, err
}

func (rcd *RedisCacheDriver) Expire(key string, expiration time.Duration) (bool, error) {
	result, err := rcd.client.Expire(context.Background(), key, expiration).Result()
	return result, err
}
