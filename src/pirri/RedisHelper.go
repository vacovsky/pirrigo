package main

import (
	"gopkg.in/redis.v5"
)

func redisWriter(message string, channel string) {
	client := redis.NewClient(&redis.Options{
		Addr:     SETTINGS.Redis.Server + ":" + SETTINGS.Redis.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	client.Publish(channel, message)
	client.Close()
}
