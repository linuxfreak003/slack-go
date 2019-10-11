package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/bluele/slack"
	log "github.com/sirupsen/logrus"
)

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: fortune | gocowsay")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	for j := 0; j < len(output); j++ {
		fmt.Printf("%c", output[j])
	}

	token := "xoxp-499911687078-498709891189-738830943330-af3669e4107f8ba4fe1661f61e0ca40e"
	message := "Hello World!"
	channelName := "public-ip"
	api := slack.New(token)
	err := api.ChatPostMessage(channelName, message, &slack.ChatPostMessageOpt{
		AsUser: true,
	})
	if err != nil {
		log.Warnf("slack message not sent: %v", err)
	}
}
