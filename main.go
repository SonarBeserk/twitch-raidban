package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

const (
	cycleInterval         = 5 * time.Second  // How often should the list be worked through
	bansPerCycle          = 10               // How many bans are allowed each cycle to stay within twitch limits
	delayBeforeDisconnect = 10 * time.Second // Delay to allow for messages to arrive before disconnecting
)

var username string
var token string
var channel string
var file string

func parseFlags() {
	flag.StringVar(&username, "username", "", "The username to authenticate as")
	flag.StringVar(&token, "token", "", "The user token for authentication")
	flag.StringVar(&channel, "channel", "", "Channel to ban bots for")
	flag.StringVar(&file, "file", "", "The txt file list of bots to ban")

	flag.Parse()

	if username == "" {
		log.Fatal("Twitch username must be provided")
	}

	if token == "" {
		log.Fatal("Twitch token must be provided")
	}

	if channel == "" {
		log.Fatal("Channel must be provided")
	}

	if file == "" {
		log.Fatal("Bots txt file must be provided")
	}
}

func main() {
	parseFlags()

	bots, err := readLines(file)
	if err != nil {
		log.Fatalf("Error parsing bots file: %v", err)
	}

	fmt.Printf("Size of list: %v", len(bots))

	banTicker := time.NewTicker(cycleInterval)

	client := twitch.NewClient(username, token)

	client.OnConnect(OnConnect(bots, banTicker, client))

	client.OnNoticeMessage(func(message twitch.NoticeMessage) {
		if message.MsgID == "already_banned" {
			fmt.Println("Notice: " + message.Message)
		}

		fmt.Println("Notice: " + message.Message)
	})

	client.Join(channel)

	err = client.Connect()
	if err != nil {
		if err.Error() == "client called Disconnect()" {
			return
		}

		fmt.Printf("Error in client: %v\n", err)
	}
}

func OnConnect(bots []string, banTicker *time.Ticker, client *twitch.Client) func() {
	return func() {
		go func() {
			totalEntries := len(bots) - 1
			entryCounter := 0
			banCounter := 0

			for {
				<-banTicker.C
				for banCounter < bansPerCycle {
					fmt.Printf("Processing entry %v/%v\n", entryCounter+1, totalEntries+1)

					if entryCounter == totalEntries {
						time.Sleep(delayBeforeDisconnect)

						fmt.Println("Disconnecting from channel")
						banTicker.Stop()
						err := client.Disconnect()
						if err != nil {
							log.Fatalln(err)
						}

						return
					}

					client.Say(channel, "/ban "+bots[entryCounter])
					banCounter++
					entryCounter++
				}

				banCounter = 0
			}
		}()
	}
}

func readLines(path string) ([]string, error) {
	var entries []string
	file, err := os.Open(path)
	if err != nil {
		return entries, fmt.Errorf("open error: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return entries, fmt.Errorf("scan error: %v", err)
	}

	return entries, nil
}
