package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rakanalh/scheduler"
	"github.com/rakanalh/scheduler/storage"
)

var sched scheduler.Scheduler

// smsRecieve handlers listens to incoming sms messages from the twilo service
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
		sendResponse(sender, "Here are your reminders:", w)
	} else {
		sendHelp(sender, w)
	}
}

// addReminder creates a cron job to send a text
func addReminder(number, message string, seconds int) {
	if _, err := sched.RunAfter(time.Duration(seconds)*time.Second, sendMessage, message, number); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sms/recieve", smsRecieve)

	// Setup storage for scheduler
	storage := storage.NewSqlite3Storage(
		storage.Sqlite3Config{
			DbName: "task_store.db",
		},
	)
	if err := storage.Connect(); err != nil {
		log.Fatal("Could not connect to db", err)
	}

	if err := storage.Initialize(); err != nil {
		log.Fatal("Could not intialize database", err)
	}

	sched = scheduler.New(storage)

	fmt.Println("Running server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
