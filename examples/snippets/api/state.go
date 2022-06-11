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

	res, status, err := pn.SetState().
		Channels([]string{"ch1"}).
		State(map[string]interface{}{
			"age": 20,
		}).
		Execute()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(res, status)
}
