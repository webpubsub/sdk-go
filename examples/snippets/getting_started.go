package main

import (
	"fmt"
	"sync"

	webpubsub "github.com/webpubsub/sdk-go/v7"
)

var pn *webpubsub.WebPubSub

func init() {
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn = webpubsub.NewWebPubSub(config)
}

func webpubsubCopy() *webpubsub.WebPubSub {
	_pn := new(webpubsub.WebPubSub)
	*_pn = *pn
	return _pn
}

func gettingStarted() {
	listener := webpubsub.NewListener()
	doneConnect := make(chan bool)
	donePublish := make(chan bool)

	msg := map[string]interface{}{
		"msg": "hello",
	}
	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case webpubsub.WPSDisconnectedCategory:
					// This event happens when radio / connectivity is lost
				case webpubsub.WPSConnectedCategory:
					// Connect event. You can do stuff like publish, and know you'll get it.
					// Or just use the connected event to confirm you are subscribed for
					// UI / internal notifications, etc
					doneConnect <- true
				case webpubsub.WPSReconnectedCategory:
					// Happens as part of our regular operation. This event happens when
					// radio / connectivity is lost, then regained.
				}
			case message := <-listener.Message:
				// Handle new message stored in message.message
				if message.Channel != "" {
					// Message has been received on channel group stored in
					// message.Channel
				} else {
					// Message has been received on channel stored in
					// message.Subscription
				}
				if msg, ok := message.Message.(map[string]interface{}); ok {
					fmt.Println("msg:=====>", msg["msg"])
				}
				/*
				   log the following items with your favorite logger
				       - message.Message
				       - message.Subscription
				       - message.Timetoken
				*/

				donePublish <- true
			case <-listener.Presence:
				// handle presence
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"hello_world"}).
		Execute()

	<-doneConnect

	response, status, err := pn.Publish().
		Channel("hello_world").Message(msg).Execute()

	if err != nil {
		// Request processing failed.
		// Handle message publish error
	}

	fmt.Println(response, status, err)

	<-donePublish
}

func listeners() {
	listener := webpubsub.NewListener()
	doneSubscribe := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case webpubsub.WPSConnectedCategory:
					doneSubscribe <- true
					return
				case webpubsub.WPSDisconnectedCategory:
					//
				case webpubsub.WPSReconnectedCategory:
					//
				case webpubsub.WPSAccessDeniedCategory:
					//
				case webpubsub.WPSUnknownCategory:
					//
				}
			case <-listener.Message:
				//
			case <-listener.Presence:
				//
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	<-doneSubscribe
}

func time() {
	res, status, err := pn.Time().Execute()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func publish() {
	res, status, err := pn.Publish().
		Channel("ch").
		Message("hey").
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func hereNow() {
	res, status, err := pn.HereNow().
		Channels([]string{"ch"}).
		IncludeUUIDs(true).
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func presence() {
	// await both connected event on emitter and join presence event received
	var wg sync.WaitGroup
	wg.Add(2)

	donePresenceConnect := make(chan bool)
	doneJoin := make(chan bool)
	doneLeave := make(chan bool)
	errChan := make(chan string)
	ch := "my-channel"

	configPresenceListener := webpubsub.NewConfig(webpubsub.GenerateUUID())
	configPresenceListener.SubscribeKey = "demo"
	configPresenceListener.PublishKey = "demo"

	pnPresenceListener := webpubsub.NewWebPubSub(configPresenceListener)

	pn.Config.UUID = "my-emitter"
	pnPresenceListener.Config.UUID = "my-listener"

	listenerEmitter := webpubsub.NewListener()
	listenerPresenceListener := webpubsub.NewListener()

	// emitter
	go func() {
		for {
			select {
			case status := <-listenerEmitter.Status:
				switch status.Category {
				case webpubsub.WPSConnectedCategory:
					wg.Done()
					return
				}
			case <-listenerEmitter.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listenerEmitter.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	// listener
	go func() {
		for {
			select {
			case status := <-listenerPresenceListener.Status:
				switch status.Category {
				case webpubsub.WPSConnectedCategory:
					donePresenceConnect <- true
				}
			case message := <-listenerPresenceListener.Message:
				errChan <- fmt.Sprintf("Unexpected message: %s",
					message.Message)
			case presence := <-listenerPresenceListener.Presence:
				fmt.Println(presence, "\n", configPresenceListener)
				// ignore join event of presence listener
				if presence.UUID == configPresenceListener.UUID {
					continue
				}

				if presence.Event == "leave" {
					doneLeave <- true
					return
				}
				wg.Done()
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	pnPresenceListener.Subscribe().
		Channels([]string{ch}).
		WithPresence(true).
		Execute()

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		panic(err)
	}

	pn.Subscribe().
		Channels([]string{ch}).
		Execute()

	go func() {
		wg.Wait()
		doneJoin <- true
	}()

	select {
	case <-doneJoin:
	case err := <-errChan:
		panic(err)
	}

	pn.Unsubscribe().
		Channels([]string{ch}).
		Execute()

	select {
	case <-doneLeave:
	case err := <-errChan:
		panic(err)
	}
}

func history() {
	res, status, err := pn.History().
		Channel("ch").
		Count(2).
		IncludeTimetoken(true).
		Reverse(true).
		Start(int64(1)).
		End(int64(2)).
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func unsubscribe() {
	pn.Subscribe().
		Channels([]string{"ch"}).
		Execute()

	// t.Sleep(3 * t.Second)

	pn.Unsubscribe().
		Channels([]string{"ch"}).
		Execute()
}

func main() {
	// gettingStarted()
	// listeners()
	// time()
	// publish()
	// hereNow()
	presence()
	// history()
	// unsubscribe()
}
