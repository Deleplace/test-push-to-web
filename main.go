package main

// Sample code from Pusher.com

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(indexPage)
	})
	http.HandleFunc("/trigger", trigger)
	http.HandleFunc("/report", report)

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
	ctx := r.Context()

	channelName := r.FormValue("chan")
	if channelName == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Expected valid chan")
		return
	}

	n, err := strconv.Atoi(r.FormValue("n"))
	if err != nil {
		log.Printf("Unexpected value for n: %q", r.FormValue("n"))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Expected valid param n")
		return
	}

	services := []serverPusher{
		pusherComPush,
		firestorePush,
	}

	log.Printf("Pushing %d events", n)

	t := time.Now()
	nerrs := 0
	for i := 0; i < n; i++ {
		// Vary the exact order of service calls
		shuffle(services)

		eventID := randomString(5)

		errs := make([]error, len(services))
		var wg sync.WaitGroup
		wg.Add(len(services))
		for i, push := range services {
			i, push := i, push
			go func() {
				defer wg.Done()
				err := push(ctx, channelName, eventID)
				if err != nil {
					errs[i] = err
					log.Println(err)
				}
			}()
		}
		wg.Wait()
		for _, err := range errs {
			if err != nil {
				nerrs++
			}
		}
	}
	duration := time.Since(t)

	log.Printf("Pushed %d events to %d services in %v => %d errors", n, len(services), duration, nerrs)
}

type serverPusher func(ctx context.Context, channelName string, eventID string) error

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
