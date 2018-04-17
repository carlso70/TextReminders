package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

// smsRecieve handlers listens to incoming sms messages from the twilo service
func smsRecieve(w http.ResponseWriter, r *http.Request) {
	sender := r.FormValue("From")
	body := r.FormValue("Body")

	// Parse message
	msg, time, err := parseMessage(body)
	if err != nil {
		sendHelp(sender, w)
	}

	if strings.Contains(body, "Set reminder:") {
		sendResponse(sender, fmt.Sprintf("Setting reminder for: %s", msg), w)
		addReminder(sender, msg, time)
	} else if strings.Contains(body, "Set timer:") {
		sendResponse(sender, fmt.Sprintf("Setting timer for: %s", msg), w)
		addReminder(sender, msg, time)
	} else if strings.Contains(body, "Get reminders") {
		sendResponse(sender, "Here are your reminders:", w)
	} else {
		sendHelp(sender, w)
	}
}

// addReminder creates a cron job to send a text
func addReminder(number, message, time string) {
	c := cron.New()

	c.AddFunc(time, func() {
		fmt.Printf("time for cron: %s\n", time)
		sendMessage(number, message)
	})
	c.Start()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sms/recieve", smsRecieve)

	log.Fatal(http.ListenAndServe(":8080", r))
}
