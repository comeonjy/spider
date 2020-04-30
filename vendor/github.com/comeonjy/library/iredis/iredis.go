package iredis

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"

	"github.com/comeonjy/library/conf"
)

var client *redis.Client
var once sync.Once



func New(c *conf.RedisConf) *redis.Client {
	once.Do(func() {
		client=ConnRedis(c)
	})
	return client
}

func ConnRedis(c *conf.RedisConf) *redis.Client {
	clien := redis.NewClient(&redis.Options{
		DB:          int(c.Database),
		Addr:        c.Address,
		Password:    c.Password,
		DialTimeout: time.Duration(c.DialTimeout) * time.Second,
	})

	if _, err := clien.Ping().Result(); err != nil {
		log.Fatal("iredis", err)
	}
	return clien
}
