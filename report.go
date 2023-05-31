package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func report(w http.ResponseWriter, r *http.Request) {
	eventID := r.FormValue("eventID")
	if eventID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Expected valid eventID")
		return
	}

	d, err := strconv.Atoi(r.FormValue("d"))
	if err != nil {
		log.Printf("Unexpected value for d: %q", r.FormValue("d"))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Expected valid param d")
		return
	}

	if d >= 0 {
		log.Printf("Firestore won, %dms faster than Pusher.com", d)
	} else {
		log.Printf("Pusher.com won, %dms faster than Firestore", -d)
	}
}
