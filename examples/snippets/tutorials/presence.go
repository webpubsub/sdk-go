package main

import (
	"fmt"

	webpubsub "github.com/webpubsub/go/v7"
)

func subscribe() {
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.UUID = "Stephen"

	pn := webpubsub.NewWebPubSub(config)

	listener := webpubsub.NewListener()
	doneConnect := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case webpubsub.WPSConnectedCategory:
					doneConnect <- true
					return
				case webpubsub.WPSReconnectedCategory:
				}
			case <-listener.Message:
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"my-channel"}).
		Execute()

	<-doneConnect
}

func hereNow() {
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.UUID = "Stephen"

	pn := webpubsub.NewWebPubSub(config)

	res, _, _ := pn.HereNow().
		Channels([]string{"ch1", "ch2"}).
		IncludeUUIDs(true).
		Execute()

	for _, v := range res.Channels {
		fmt.Println("Channel: ", v.ChannelName)
		fmt.Println("Occupancy: ", v.Occupancy)
		fmt.Println("Occupants")

		for _, v := range v.Occupants {
			fmt.Println("UUID: ", v.UUID, ", state: ", v.State)
		}
	}
}

func globalHereNow() {
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.UUID = "Stephen"

	pn := webpubsub.NewWebPubSub(config)

	res, status, err := pn.HereNow().
		IncludeState(true).
		IncludeUUIDs(true).
		Execute()

	fmt.Println(res, status, err)
}

func whereNow() {
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.UUID = "Stephen"

	pn := webpubsub.NewWebPubSub(config)

	res, status, err := pn.WhereNow().
		UUID("person-uuid").
		Execute()

	fmt.Println(res, status, err)
}

func getState() {
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.UUID = "Stephen"

	pn := webpubsub.NewWebPubSub(config)

	res, status, err := pn.GetState().
		Channels([]string{"my-unique-ch"}).
		Execute()

	fmt.Println(res, status, err)
}

func main() {
	subscribe()
	hereNow()
	globalHereNow()
	getState()
}
