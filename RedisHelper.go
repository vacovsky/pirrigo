package main

import (
	"strconv"

	"gopkg.in/redis.v5"
)

func RedisWriter(message string, channel string) {
	client := redis.NewClient(&redis.Options{
		Addr:     CONFIG.RedisServer + ":" + strconv.Itoa(CONFIG.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	client.Publish(channel, message)
	client.Close()
}
