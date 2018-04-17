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

func smsRecieve(w http.ResponseWriter, r *http.Request) {
	sender := r.FormValue("From")
	body := r.FormValue("Body")

	resp := twiml.NewResponse()
	resp.Action(twiml.Message{
		Body: fmt.Sprintf("Hello, %s, you said: %s", sender, body),
		From: TWILO_NUM,
		To:   sender,
	})
	resp.Send(w)
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there human!")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sms/recieve", smsRecieve)
	r.HandleFunc("/", test)

	log.Fatal(http.ListenAndServe(":80", r))
}
