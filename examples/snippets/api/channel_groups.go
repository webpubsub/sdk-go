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

	addRes, status, err := pn.AddChannelToChannelGroup().
		Channels([]string{"ch1", "ch2"}).
		ChannelGroup("cg1").
		Execute()

	if err != nil {
		fmt.Println("Error :", err)
	}

	fmt.Println(addRes, status)

	listRes, status, err := pn.ListChannelsInChannelGroup().
		ChannelGroup("cg1").
		Execute()

	if err != nil {
		fmt.Println("Error :", err)
	}

	fmt.Println(listRes, status)

	removeRes, status, err := pn.RemoveChannelFromChannelGroup().
		Channels([]string{"ch1", "ch2"}).
		ChannelGroup("cg1").
		Execute()

	if err != nil {
		fmt.Println("Error :", err)
	}

	fmt.Println(removeRes, status)
}
