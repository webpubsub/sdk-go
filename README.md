# WebPubSub Go SDK

[![GoDoc](https://godoc.org/github.com/webpubsub/go?status.svg)](https://godoc.org/github.com/webpubsub/go)
[![Build Status](https://travis-ci.com/webpubsub/go.svg?branch=master)](https://travis-ci.com/webpubsub/go)
[![codecov.io](https://codecov.io/github/webpubsub/go/coverage.svg)](https://codecov.io/github/webpubsub/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/webpubsub/go)](https://goreportcard.com/report/github.com/webpubsub/go)

This is the official WebPubSub Go SDK repository.

WebPubSub takes care of the infrastructure and APIs needed for the realtime communication layer of your application. Work on your app's logic and let WebPubSub handle sending and receiving data across the world in less than 100ms.

## Requirements

* Go (1.9+)

## Get keys

You will need the publish and subscribe keys to authenticate your app. Get your keys from the [Admin Portal](https://dashboard.webpubsub.com/).

## Configure WebPubSub

1. Integrate WebPubSub into your project using the `go` command:

    ```go
    go get github.com/webpubsub/go
    ```

    If you encounter dependency issues, use the `dep ensure` command to resolve them.

1. Create a new file and add the following code:

    ```go
    func main() {
        config := webpubsub.NewConfig()
        config.SubscribeKey = "mySubscribeKey"
        config.PublishKey = "myPublishKey"
        config.UUID = "myUniqueUUID"

        pn := webpubsub.NewWebPubSub(config)
    }
    ```

## Add event listeners

```go
listener := webpubsub.NewListener()
doneConnect := make(chan bool)
donePublish := make(chan bool)

msg := map[string]interface{}{
    "msg": "Hello world",
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
                fmt.Println(msg["msg"])
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
```

## Publish and subscribe

In this code, publishing a message is triggered when the subscribe call is finished successfully. The `Publish()` method uses the `msg` variable that is used in the following code.

```go
msg := map[string]interface{}{
        "msg": "Hello world!"
}

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
```

## Documentation

[API reference for Go](https://www.webpubsub.com/docs/go/webpubsub-go-sdk-v4)

## Support

If you **need help** or have a **general question**, contact <support@webpubsub.com>.
