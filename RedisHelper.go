package main

import (
	"gopkg.in/redis.v5"
)

func RedisWriter(message string, channel string) {
	client := redis.NewClient(&redis.Options{
		Addr:     REDISHOST + ":" + REDISPORT,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	client.Publish(channel, message)
	client.Close()
}
