package main

// Contains all the code relating to twilo sms sending

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

		formating for timed reminder = second/minutes/hour/day
		formating for reminder time = second/minute/hour/day`
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

// SendMessage sends a basic text
func sendMessage(numberTo, message string) {
	fmt.Printf("Sending message to %s, '%s'\n", numberTo, message)
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
	fmt.Printf("Sending response to %s, '%s'\n", numberTo, message)
	resp := twiml.NewResponse()
	resp.Action(twiml.Message{
		Body: message,
		From: twiloNum,
		To:   numberTo,
	})
	resp.Send(w)
}

// parseMessage takes in a sms, and returns the message, and time in seconds
func parseMessage(sms string) (message string, secs time.Duration, err error) {
	params := strings.Split(sms, ":")
	if len(params) != 3 {
		return "", 0, errors.New("Invalid formatting")
	}

	// format of time msg second/minute/hour/day/year
	times := strings.Split(params[2], "/")
	for i := 0; i < len(times) || i > 3; i++ {
		switch i {
		case 0:
			// Seconds
			s, err := strconv.Atoi(strings.TrimSpace(times[i]))
			if err != nil {
				panic(err)
			}
			secs += time.Duration(s) * time.Second
		case 1:
			// Minutes
			s, err := strconv.Atoi(strings.TrimSpace(times[i]))
			if err != nil {
				panic(err)
			}
			secs += time.Duration(s) * time.Minute
		case 2:
			// Hours
			s, err := strconv.Atoi(strings.TrimSpace(times[i]))
			if err != nil {
				panic(err)
			}
			secs += time.Duration(s) * time.Hour
		case 3:
			// Days
			s, err := strconv.Atoi(strings.TrimSpace(times[i]))
			if err != nil {
				panic(err)
			}
			secs += time.Duration(s*24) * time.Hour
		}
	}

	return params[1], secs, nil
}
