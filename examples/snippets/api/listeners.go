package main

import (
	"fmt"

	webpubsub "github.com/webpubsub/go/v7"
)

func main() {
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn := webpubsub.NewWebPubSub(config)

	listener := webpubsub.NewListener()
	done := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case webpubsub.WPSConnectedCategory:
					done <- true
				}
			case message := <-listener.Message:
				//Channel
				fmt.Println(message.Channel)
				//Subscription
				fmt.Println(message.Subscription)
				//Payload
				fmt.Println(message.Message)
				//Publisher ID
				fmt.Println(message.Publisher)
				//Timetoken
				fmt.Println(message.Timetoken)
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	<-done
}
