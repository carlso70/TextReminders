package main

import (
	"errors"
	"fmt"
	"strings"
)

// parseMessage takes in a sms, and returns the message, and time formatted for use with cron package
func parseMessage(sms string) (message, time string, err error) {
	fmt.Println(sms)
	params := strings.Split(sms, ":")
	if len(params) != 3 {
		return "", "", errors.New("Invalid formatting")
	}

	return params[1], params[2], nil
}
