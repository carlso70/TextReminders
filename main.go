package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"bitbucket.org/ckvist/twilio/twiml"
	"bitbucket.org/ckvist/twilio/twirest"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
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

	_, err := twilioClient.Request(msg)
	if err != nil {
		fmt.Errorf("Error sending sms text: %v", err)
	}

}

// sendResponse Sends a response from the twilo number associated with the account
func sendResponse(number_to, message string, w http.ResponseWriter) {
	resp := twiml.NewResponse()
	resp.Action(twiml.Message{
		Body: message,
		From: TWILO_NUM,
		To:   number_to,
	})
	resp.Send(w)
}

// smsRecieve handlers listens to incoming sms messages from the twilo service
func smsRecieve(w http.ResponseWriter, r *http.Request) {
	sender := r.FormValue("From")
	body := r.FormValue("Body")

	fmt.Println("Recieved message from %s: %s", sender, body)

	if strings.Contains(body, "set reminder:") {
		sendResponse(sender, "Setting reminder for: TBD", w)
		addReminder(sender, body)
	} else if strings.Contains(body, "get reminders") {
		sendResponse(sender, "Here are your reminders:", w)
	} else {
		sendResponse(sender, "INVALID text", w)
	}
}

// addReminder creates a cron job to send a text
func addReminder(number, message string) {
	c := cron.New()

	c.AddFunc("10 * * * * *", func() {
		fmt.Println("Every 10 seconds, sending text")
		sendMessage(number, message)
	})
	c.Start()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sms/recieve", smsRecieve)

	log.Fatal(http.ListenAndServe(":8080", r))
}
