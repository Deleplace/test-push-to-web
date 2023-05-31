package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func report(w http.ResponseWriter, r *http.Request) {
	channelName := r.FormValue("chan")
	if channelName == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Expected valid chan")
		return
	}

	eventID := r.FormValue("eventId")
	if eventID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Expected valid eventId")
		return
	}

	d, err := strconv.Atoi(r.FormValue("d"))
	if err != nil {
		log.Printf("Unexpected value for d: %q", r.FormValue("d"))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Expected valid param d")
		return
	}
	log.Printf("Chan %q, Event %q, diff=%d\n", channelName, eventID, d)

	if d >= 0 {
		log.Printf("Firestore won, %dms faster than Pusher.com\n", d)
	} else {
		log.Printf("Pusher.com won, %dms faster than Firestore\n", -d)
	}
}
