package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var (
	TWILO_KEY     = os.Getenv("TWILO_KEY")
	TWILO_ACCOUNT = os.Getenv("TWILO_ACCOUNT")
	TWILO_NUM     = os.Getenv("TWILO_NUM")
	URL_STR       = strings.Join([]string{"https://api.twilio.com/2010-04-01/Accounts/", TWILO_ACCOUNT, "/Messages.json"}, "")
)

func sendMessage(number_to, message string) {
	msgData := url.Values{}
	msgData.Set("To", number_to)
	msgData.Set("From", TWILO_NUM)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", URL_STR, &msgDataReader)
	req.SetBasicAuth(TWILO_ACCOUNT, TWILO_KEY)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
}

func smsRecieve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	text := vars["text"]

	fmt.Fprintf(w, "Recieved the text: %s\n", text)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/smsRecieve", smsRecieve).Methods("GET")

	http.ListenAndServe(":80", r)

	fmt.Println("Running Server....")

	// send test message
	sendMessage("+18477215493", "test message")
}
