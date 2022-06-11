package main

import webpubsub "github.com/webpubsub/sdk-go/v7"

func main() {
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.UUID = "my_uuid"

	pn := webpubsub.NewWebPubSub(config)

	meta := map[string]interface{}{
		"my":   "meta",
		"name": "WebPubSub",
	}

	pn.Subscribe().
		Channels([]string{"ch1"}).
		Execute()

	pn.Publish().
		Meta(meta).
		Message("hello").
		Channel("ch1").
		Execute()
}
