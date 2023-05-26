package main

// Sample code from Pusher.com

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/pusher/pusher-http-go/v5"
)

var pusherComSecret = os.Getenv("PUSHER_COM_SECRET")

func main() {
	if pusherComSecret == "" {
		log.Fatal("Could not find env var PUSHER_COM_SECRET")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(indexPage)
	})
	http.HandleFunc("/trigger", trigger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatal(err)
}

//go:embed index.html
var indexPage []byte

func trigger(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(r.FormValue("n"))
	if err != nil {
		log.Printf("Unexpected value for n: %q", r.FormValue("n"))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Expected valid param n")
		return
	}

	services := []serverPusher{
		pusherPush,
	}

	for i := 0; i < n; i++ {
		// Vary the exact order of service calls
		shuffle(services)

		eventID := randomString(5)

		for _, push := range services {
			err := push(eventID)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

type serverPusher func(eventID string) error

func pusherPush(eventID string) error {
	pusherClient := pusher.Client{
		AppID:   "1608025",
		Key:     "93f4a1f9e72133245d66",
		Secret:  pusherComSecret,
		Cluster: "eu",
		Secure:  true,
	}

	// The payload doesn't matter, only the event delivery
	data := map[string]string{"foo": "bar"}
	return pusherClient.Trigger("new-data", eventID, data)
}

func shuffle[T any](a []T) {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}

const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func randomString(n int) string {
	a := make([]byte, n)
	for i := range a {
		a[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return string(a)
}
