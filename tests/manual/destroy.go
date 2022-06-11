package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	webpubsub "github.com/webpubsub/sdk-go/v7"
)

func main() {
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	pn := webpubsub.NewWebPubSub(config)

	fmt.Println("vim-go")

	// Add listeners
	ln1 := webpubsub.NewListener()
	ln2 := webpubsub.NewListener()
	ln3 := webpubsub.NewListener()

	// TODO: listen on 1st listener

	pn.AddListener(ln1)

	pn.Subscribe().
		Channels([]string{"blah"}).
		Execute()

	pn.AddListener(ln2)
	pn.AddListener(ln3)

	time.Sleep(1 * time.Second)

	pn.Destroy()
	time.Sleep(1 * time.Second)

	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}
