package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	token := os.Getenv("SLACK_API_TOKEN")
	api := slack.New(token)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	if _, ok := os.LookupEnv("DEBUG"); ok {
		api.SetDebug(true)
	}

	resp, err := api.GetBotInfo("")
	if err != nil {
		log.Fatalf("Error GetUserIdentity %s", err.Error())
		os.Exit(1)
	}
	botID := resp.ID

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if ev.User == botID {
				continue
			}
			_, _, err := api.PostMessage("C048MG6B6", ev.Text, slack.PostMessageParameters{})
			if err != nil {
				log.Fatalf("Error %s", err.Error())
			}
		}
	}
}
