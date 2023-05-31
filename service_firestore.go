package main

import (
	"context"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
)

var gcpProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

func init() {
	if gcpProjectID == "" {
		log.Fatal("Could not find env var GOOGLE_CLOUD_PROJECT")
	}
}

type ChannelModel struct{}

type MessageModel struct {
	Message string
}

func firestorePush(ctx context.Context, channelName string, eventID string) error {
	client, err := firestore.NewClient(ctx, gcpProjectID)
	if err != nil {
		return err
	}

	chanDoc := client.Doc("channels/" + channelName)

	_, err = chanDoc.Create(ctx, ChannelModel{})
	if err != nil {
		if strings.Contains(err.Error(), "Document already exists") {
			log.Println(err)
			// it's okay if the doc already existed, carry on
		} else {
			return err
		}
	}

	// the "id" of a message could be just a fine-grained timestamp
	// t := time.Now().UnixNano()

	msgDoc := client.Doc("channels/" + channelName + "/messages/" + eventID)

	_, err = msgDoc.Create(ctx, MessageModel{
		Message: "Hello from the backend",
	})
	if err != nil {
		log.Println(err)
		// it's okay-ish if the message already existed, carry on
	}

	return nil
}
