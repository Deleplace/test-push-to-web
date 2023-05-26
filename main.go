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
)

func main() {
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
		pusherComPush,
	}

	log.Printf("Pushing %d events", n)

	nerrs := 0
	for i := 0; i < n; i++ {
		// Vary the exact order of service calls
		shuffle(services)

		eventID := randomString(5)

		for _, push := range services {
			err := push(eventID)
			if err != nil {
				nerrs++
				log.Println(err)
			}
		}
	}
	log.Printf("Pushed %d events to %d services => %d errors", n, len(services), nerrs)
}

type serverPusher func(eventID string) error

const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func randomString(n int) string {
	a := make([]byte, n)
	for i := range a {
		a[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return string(a)
}

func shuffle[T any](a []T) {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}
