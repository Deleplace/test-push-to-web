package main

// Sample code from Pusher.com

import (
	"fmt"

	"github.com/pusher/pusher-http-go/v5"
)

func main() {
	pusherClient := pusher.Client{
		AppID:   "1608025",
		Key:     "497269378f3faf6ec5e7",
		Secret:  "73c9735bfa358ae25d12",
		Cluster: "eu",
		Secure:  true,
	}

	data := map[string]string{"message": "hello world"}
	err := pusherClient.Trigger("my-channel", "my-event", data)
	if err != nil {
		fmt.Println(err.Error())
	}
}
