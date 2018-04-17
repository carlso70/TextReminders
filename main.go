package main

import (
	"cron"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"bitbucket.org/ckvist/twilio/twiml"
	"bitbucket.org/ckvist/twilio/twirest"
	"github.com/gorilla/mux"
)

var (
	TWILO_KEY     = os.Getenv("TWILO_KEY")
	TWILO_ACCOUNT = os.Getenv("TWILO_ACCOUNT")
	TWILO_NUM     = os.Getenv("TWILO_NUM")
	URL_STR       = strings.Join([]string{"https://api.twilio.com/2010-04-01/Accounts/", TWILO_ACCOUNT, "/Messages.json"}, "")
	twilioClient  = twirest.NewClient(
		TWILO_ACCOUNT,
		TWILO_KEY)
)

func sendMessage(number_to, message string) {
	msg := twirest.SendMessage{
		Text: message,
		From: TWILO_NUM,
		To:   number_to,
	}

	twilioClient.Request(msg)
}

// sendResponse Sends a response from the twilo number associated with the account
func sendResponse(number_to, message string, w http.ResponseWriter) {
	resp := twiml.NewResponse()
	resp.Action(twiml.Message{
		Body: fmt.Sprintf("Hello, %s, you said: %s", sender, body),
		From: TWILO_NUM,
		To:   sender,
	})
	resp.Send(w)
}

// smsRecieve handlers listens to incoming sms messages from the twilo service
func smsRecieve(w http.ResponseWriter, r *http.Request) {
	sender := r.FormValue("From")
	body := r.FormValue("Body")

	fmt.Println("Recieved message from %s: %s", sender, body)

	switch body {
	case strings.Contains(body, "set reminder:"):
		sendResponse(sender, "Setting reminder for: TBD", w)
	case strings.Contains(body, "get reminders"):
		sendResponse(sender, "Here are your reminders:", w)
	default:
		sendResponse(sender, "INVALID text")
	}

}

// addReminder creates a cron job to send a text
func addReminder() {

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sms/recieve", smsRecieve)

	log.Fatal(http.ListenAndServe(":8080", r))
}
