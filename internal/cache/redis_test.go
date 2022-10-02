package cache

import (
	"context"
	"log"
	"testing"
	"time"
)

var ctx = context.Background()

func TestSetString(t *testing.T) {
	err := GetRedisClient().Set(ctx, "test", "foo", 0).Err()
	if err != nil {
		panic(err)
	}
}

func TestGetString(t *testing.T) {
	val, err := GetRedisClient().Get(ctx, "test").Result()
	if err != nil {
		panic(err)
	}
	log.Fatalf("value:%v\n", val)
}

func TestSAdd(t *testing.T) {
	err := GetRedisClient().SAdd(ctx, "set", 2, 1).Err()
	if err != nil {
		panic(err)
	}
	err = GetRedisClient().Expire(ctx, "set", 100*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

// get Set keys and values
func TestSMember(t *testing.T) {
	ms, err := GetRedisClient().SMembers(ctx, "set").Result()
	if err != nil {
		panic(err)
	}
	for _, m := range ms {
		t.Logf("TestSMember value: %v\n", m)
	}
}

func TestHSet(t *testing.T) {
	err := GetRedisClient().HSet(ctx, "test001", "action", "add").Err()
	if err != nil {
		panic(err)
	}
}

func TestHGet(t *testing.T) {
	action, err := GetRedisClient().HGet(ctx, "test001", "action").Result()
	if err != nil {
		panic(err)
	}
	log.Fatalf("action:%v", action)
}
