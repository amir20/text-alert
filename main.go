package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

var (
	accountSid      = ""
	authToken       = ""
	to              = ""
	from            = ""
	autoDetectMedia = false
	version         = "dev"
	commit          = "none"
	date            = "unknown"
)

func init() {
	flag.StringVar(&accountSid, "accountSid", "", "Twilio accountSid")
	flag.StringVar(&authToken, "authToken", "", "Twilio authToken")
	flag.StringVar(&to, "to", "", "To phone number")
	flag.StringVar(&from, "from", "", "From phone number")
	flag.BoolVarP(&autoDetectMedia, "detect-media", "m", false, "If enabled will send image URLs as MMS instead")
	flag.Parse()
}

func main() {
	log.SetFlags(0)
	requireFlags("accountSid", "authToken", "to", "from")

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	client := &http.Client{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		msgData := url.Values{}
		msgData.Set("To", to)
		msgData.Set("From", from)
		text := scanner.Text()

		parsed, err := url.Parse(text)
		if autoDetectMedia && err == nil && isImage(parsed) {
			msgData.Set("MediaUrl", parsed.String())
		} else {
			msgData.Set("Body", text)
		}

		msgDataReader := strings.NewReader(msgData.Encode())

		req, _ := http.NewRequest("POST", urlStr, msgDataReader)
		req.SetBasicAuth(accountSid, authToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, _ := client.Do(req)
		if resp.StatusCode != 201 {
			var data map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&data)
			if err == nil {
				log.Fatalln(data)
			}

		}
	}
}

func isImage(url *url.URL) bool {
	ext := filepath.Ext(url.Path)
	return url.IsAbs() &&
		(ext == ".jpg" || ext == ".png" || ext == ".gif")
}

func requireFlags(names ...string) {
	for _, name := range names {
		f := flag.Lookup(name)
		if !f.Changed {
			log.Fatalf("Flag --%s is required", name)
		}
	}

}
