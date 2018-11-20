package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

func test() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println("-> " + scanner.Text())
	}
}

var (
	accountSid = ""
	authToken  = ""
	to         = ""
	from       = ""
	version    = "dev"
	commit     = "none"
	date       = "unknown"
)

func init() {
	flag.StringVar(&accountSid, "accountSid", "", "Twilio accountSid")
	flag.StringVar(&authToken, "authToken", "", "Twilio authToken")
	flag.StringVar(&to, "to", "", "To phone number")
	flag.StringVar(&from, "from", "", "From phone number")

	flag.Parse()
}

func main() {
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	client := &http.Client{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		msgData := url.Values{}
		msgData.Set("To", to)
		msgData.Set("From", from)
		msgData.Set("Body", scanner.Text())
		msgDataReader := strings.NewReader(msgData.Encode())

		req, _ := http.NewRequest("POST", urlStr, msgDataReader)
		req.SetBasicAuth(accountSid, authToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, _ := client.Do(req)
		if resp.StatusCode != 201 {
			log.Fatalln(resp.Body)
		}
	}
}
