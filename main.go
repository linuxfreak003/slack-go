package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

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

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		fmt.Println("Error: No message given")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	var build strings.Builder
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		build.WriteRune(input)
	}
	message = build.String()
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
