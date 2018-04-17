package main

// Contains all the code relating to twilo sms sending

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

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

// parseMessage takes in a sms, and returns the message, and time formatted for use with cron package
func parseMessage(sms string) (message, time string, err error) {
	params := strings.Split(sms, ":")
	if len(params) != 3 {
		return "", "", errors.New("Invalid formatting")
	}

	// format of time msg second/minute/hour/day/year
	times := strings.Split(params[2], "/")

	// There are 6 mandatory fields for the Cron job library param
	var buffer bytes.Buffer
	for i := 0; i < 6; i++ {
		if i < len(times)-1 {
			buffer.WriteString(strings.TrimSpace(times[i]))
			buffer.WriteString(" ")
		} else {
			buffer.WriteString("* ")
		}
	}

	return params[1], buffer.String(), nil
}
