package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// smsRecieve handlers listens to incoming sms messages from the twilio service
func smsRecieve(w http.ResponseWriter, r *http.Request) {
	sender := r.FormValue("From")
	body := r.FormValue("Body")

	fmt.Printf("Recieved message from %s, '%s'\n", sender, body)

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
		msg = "Here are you reminders:\n"
		for _, v := range activeReminders[sender] {
			msg += v
			msg += "\n"
		}
		sendResponse(sender, msg, w)
	} else {
		sendHelp(sender, w)
	}
}

// fallback handler for when there is an error from twilio sms service
func fallback(w http.ResponseWriter, r *http.Request) {
	errCode := r.FormValue("ErrorCode")
	errURL := r.FormValue("ErrorUrl")

	fmt.Printf("Twilio Error code %s: %s\n", errCode, errURL)
}

// addReminder runs a ticker that sends a sms at the end
func addReminder(number, message string, length time.Duration) {
	activeReminders[number] = append(activeReminders[number], message)

	ticker := time.NewTicker(length)
	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				sendMessage(number, message)
				activeReminders[number] = removeString(activeReminders[number], message)
				ticker.Stop()
			}
		}
	}(ticker)
}

// removeString removes an item from an array of strings
func removeString(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sms/recieve", smsRecieve)
	r.HandleFunc("/sms/fallback", fallback)

	fmt.Println("Running server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
