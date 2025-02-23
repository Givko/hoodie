package subscribers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/redis/go-redis/v9"
)

func StartChatMessageSubscriber(
	redisClient *redis.Client,
	logger logr.Logger,
) {
	ctx := context.Background()
	streamKey := "chat_messages"
	groupName := "persisters"
	err := redisClient.XGroupCreateMkStream(ctx, streamKey, groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		logger.Error(err, "Error creating consumer group")
	}

	// Infinite loop to continuously poll for messages.
	for {
		// XReadGroup waits up to 30 seconds for new messages; up to 10 messages are returned if available.
		res, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: "consumer_persister",
			Streams:  []string{streamKey, ">"},
			Count:    10,
			Block:    30 * time.Second,
		}).Result()

		if err != nil {
			// If no messages are available, Redis returns a redis.Nil error.
			if err == redis.Nil {
				continue
			}

			logger.Error(err, "Error reading messages")
			continue
		}

		// Process the returned messages.
		for _, stream := range res {
			for _, message := range stream.Messages {
				fmt.Printf("Processing message ID: %s\n", message.ID)
				for key, value := range message.Values {
					fmt.Printf("  %s: %v\n", key, value)
				}

				// Acknowledge the message to prevent redelivery.
				if err := redisClient.XAck(ctx, streamKey, groupName, message.ID).Err(); err != nil {
					logger.Info("Error acknowledging message", "id", message.ID, "error", err.Error())
				}
			}
		}
	}
}
