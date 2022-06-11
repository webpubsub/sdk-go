package main

import (
	"fmt"

	webpubsub "github.com/webpubsub/sdk-go/v7"
)

var pn *webpubsub.WebPubSub

func init() {
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn = webpubsub.NewWebPubSub(config)
}

func main() {
	listener := webpubsub.NewListener()
	doneSubscribe := make(chan bool)
	data := make(map[string]interface{})

	data["awesome"] = "data"

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case webpubsub.WPSConnectedCategory:
					res, status, err := pn.Publish().
						Channel("awesome-channel").
						Message(data).
						Execute()

					fmt.Printf(res, status, err)

					doneSubscribe <- true
					return
				}
			case <-listener.Message:
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"awesome-channel"}).
		Execute()

	<-doneSubscribe

	pn.Unsubscribe().
		Channels([]string{"awesome-channel"}).
		Execute()
}
