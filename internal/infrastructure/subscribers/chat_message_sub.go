package subscribers

import (
	"context"

	wsProto "github.com/givko/hoodie/internal/api/ws/proto"
	"github.com/go-logr/logr"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

func StartChatMessageSubscriber(
	redisClient *redis.Client,
	logger logr.Logger,
) {
	ctx := context.Background()
	sub := redisClient.Subscribe(ctx, "chat_messages")
	defer sub.Close()

	messagesChannel := sub.Channel()
	for message := range messagesChannel {
		protoMessage := &wsProto.Message{}
		messageAsByteArray := []byte(message.Payload)
		err := proto.Unmarshal(messageAsByteArray, protoMessage)
		if err != nil {
			logger.Error(err, "Error unmarshalling message from redis", "message", message)
			continue
		}

		logger.Info("Received message from redis", "consumedMessage", protoMessage)
	}
}
