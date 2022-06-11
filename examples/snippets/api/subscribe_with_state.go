package main

import (
	"fmt"

	webpubsub "github.com/webpubsub/sdk-go/v7"
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
			case <-listener.Message:
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	res, status, err := pn.SetState().
		Channels([]string{"ch"}).
		State(map[string]interface{}{
			"field_a": "cool",
			"field_b": 21,
		}).
		Execute()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(res, status)

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	<-done
}
