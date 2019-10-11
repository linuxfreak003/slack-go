package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bluele/slack"
	log "github.com/sirupsen/logrus"
)

func main() {
	var message, token, channel string
	var slackToken string

	home, _ := os.UserHomeDir()
	tokenFile, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/slack-go/slack-token", home))
	if err == nil {
		slackToken = string(tokenFile)
	}

	if slackToken == "" {
		slackToken = os.Getenv("SLACK_TOKEN")
	}

	slackToken = strings.TrimSpace(slackToken)

	flag.StringVar(&message, "m", "", "Message to send")
	flag.StringVar(&token, "t", slackToken, "Slack API Token")
	flag.StringVar(&channel, "c", "", "Slack Channel to send to")
	flag.Parse()

	if channel == "" {
		fmt.Println("Error: No channel specified")
		return
	}

	if message != "" {
		api := slack.New(token)
		err = api.ChatPostMessage(channel, message, &slack.ChatPostMessageOpt{
			AsUser: true,
		})
		if err != nil {
			log.Warnf("slack message not sent: %v", err)
		}
		return
	}

	fi, _ := os.Stdin.Stat() // get the FileInfo struct describing the standard input.

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("data is from pipe")
		// do things for data from pipe

		bytes, _ := ioutil.ReadAll(os.Stdin)
		message = string(bytes)
	} else {
		fmt.Println("data is from terminal")
		// do things from data from terminal

		ConsoleReader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter message to send : ")

		message, err = ConsoleReader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Message Sent.")
	}

	if message == "" {
		fmt.Println("Error: No message given")
	}

	api := slack.New(token)
	err = api.ChatPostMessage(channel, message, &slack.ChatPostMessageOpt{
		AsUser: true,
	})
	if err != nil {
		log.Warnf("slack message not sent: %v", err)
	}
}
