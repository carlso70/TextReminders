package main

// Contains all the code relating to twilo sms sending

import (
	"log"
	"net/http"
	"os"

	"bitbucket.org/ckvist/twilio/twiml"
	"bitbucket.org/ckvist/twilio/twirest"
)

var (
	twiloKey     = os.Getenv("TWILO_KEY")
	twiloAccount = os.Getenv("TWILO_ACCOUNT")
	twiloNum     = os.Getenv("TWILO_NUM")
	twilioClient = twirest.NewClient(
		twiloAccount,
		twiloKey)
)

const (
	// HELP is the default help message whenever a user types an invalid message
	HELP = `INVALID Format: 
		Examples
		Set reminder: clean : 10/1/5/1
		Set timer: time over : 15/0/0 
		Get reminders

		formating for timed reminder = second/minutes/hour
		formating for reminder time = second/minute/hour/day/year`
)

// sendHelp responds to a message send to the twilo service, which is tied to the dev API number
func sendHelp(numberTo string, w http.ResponseWriter) {
	resp := twiml.NewResponse()
	resp.Action(twiml.Message{
		Body: HELP,
		From: twiloNum,
		To:   numberTo,
	})
	resp.Send(w)
}

// sendMessage sends a basic text
func sendMessage(numberTo, message string) {
	msg := twirest.SendMessage{
		Text: message,
		From: twiloNum,
		To:   numberTo,
	}

	_, err := twilioClient.Request(msg)
	if err != nil {
		log.Fatal(err)
	}

}

// sendResponse Sends a response from the twilo number associated with the account
func sendResponse(numberTo, message string, w http.ResponseWriter) {
	resp := twiml.NewResponse()
	resp.Action(twiml.Message{
		Body: message,
		From: twiloNum,
		To:   numberTo,
	})
	resp.Send(w)
}
