package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nu50218/miopon-go"
)

type Settings struct {
	DeveloperID    string   `json:"developer_id"`
	AccessTokens   []string `json:"access_tokens"`
	IntervalMinute int64    `json:"interval_minute"`
}

func main() {
	var isInit bool
	flag.BoolVar(&isInit, "init", false, "Make settings.json in the working directory")
	var developerID string
	flag.StringVar(&developerID, "id", "", "developerID")
	flag.Parse()

	if isInit {
		if developerID == "" {
			fmt.Print("developerID:")
			fmt.Scan(&developerID)
		}

		var redirectURI, state string
		fmt.Print("redirectURI:")
		fmt.Scan(&redirectURI)
		fmt.Print("state:")
		fmt.Scan(&state)

		settings := &Settings{
			DeveloperID:    developerID,
			AccessTokens:   []string{},
			IntervalMinute: 1,
		}
		b, err := json.MarshalIndent(settings, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		if err := ioutil.WriteFile("settings.json", b, 0644); err != nil {
			log.Fatalln(err)
		}

		client := miopon.New(developerID, &miopon.Settings{})
		fmt.Println("Created settings file. Get access tokens from the following URL and add them to the settings file.")
		fmt.Println(client.MakeAuthorizationURL(redirectURI, state))
		return
	}

}
