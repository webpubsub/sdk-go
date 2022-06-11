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
	channelGroup := "family"

	res, status, err := pn.ListChannelsInChannelGroup().
		ChannelGroup("family").
		Execute()

	fmt.Println(res, status, err)

	resAdd, statusAdd, err := pn.AddChannelToChannelGroup().
		Channels([]string{"wife"}).
		ChannelGroup(channelGroup).
		Execute()

	fmt.Println(resAdd, statusAdd, err)

	resAdd, statusAdd, err = pn.AddChannelToChannelGroup().
		Channels([]string{"son", "daughter"}).
		ChannelGroup(channelGroup).
		Execute()

	fmt.Println(resAdd, statusAdd, err)

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
		ChannelGroups([]string{channelGroup}).
		Timetoken(int64(1337)).
		WithPresence(true).
		Execute()

	<-doneConnect

	pn.Unsubscribe().
		ChannelGroups([]string{channelGroup}).
		Execute()

	pn.Subscribe().
		ChannelGroups([]string{"cg1", "cg2"}).
		Timetoken(int64(1337)).
		WithPresence(true).
		Execute()

	resRemove, statusRemove, err := pn.RemoveChannelFromChannelGroup().
		Channels([]string{"son"}).
		ChannelGroup("family").
		Execute()

	fmt.Println(resRemove, statusRemove, err)

	res, status, err = pn.ListChannelsInChannelGroup().
		ChannelGroup("family").
		Execute()

	fmt.Println(res, status, err)

	resDel, statusDel, err := pn.DeleteChannelGroup().
		ChannelGroup("family").
		Execute()

	fmt.Println(resDel, statusDel, err)
}
