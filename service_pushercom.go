package main

import (
	"log"
	"os"

	"github.com/pusher/pusher-http-go/v5"
)

var pusherComSecret = os.Getenv("PUSHER_COM_SECRET")

func init() {
	if pusherComSecret == "" {
		log.Fatal("Could not find env var PUSHER_COM_SECRET")
	}
}

func pusherComPush(eventID string) error {
	pusherClient := pusher.Client{
		AppID:   "1608025",
		Key:     "93f4a1f9e72133245d66",
		Secret:  pusherComSecret,
		Cluster: "eu",
		Secure:  true,
	}

	const (
		channelName = "server-push-test-channel"
		eventName   = "new-data"
	)
	data := map[string]string{"id": eventID}
	return pusherClient.Trigger(channelName, eventName, data)
}
