package e2e

import (
	"context"
	//"fmt"
	"log"
	"os"
	"runtime"
	"testing"

	//"time"

	webpubsub "github.com/webpubsub/go/v7"
)

// import _ "net/http/pprof"
// import "net/http"

var NumOfPublishes = 10

func TestDestroy(t *testing.T) {
	testSerial := "aa"
	config := configCopy()
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	if enableDebuggingInTests {
		config.Log = log.New(os.Stdout, "", log.Ltime|log.Lmicroseconds)
	}
	config.WPSReconnectionPolicy = webpubsub.WPSExponentialPolicycy
	config.MaximumReconnectionRetries = -1
	config.UUID = testSerial
	config.SuppressLeaveEvents = true
	config.SetPresenceTimeoutWithCustomInterval(330, 300)
	config.MaxWorkers = 0

	pn := webpubsub.NewWebPubSub(config)
	listener := webpubsub.NewListener()
	pn.AddListener(listener)
	pn.Subscribe().Channels([]string{"a." + testSerial}).Execute()
	<-listener.Status
	ctx, cancelPub := context.WithCancel(context.Background())
	for i := 0; i < NumOfPublishes; i++ {
		go pn.PublishWithContext(ctx).Channel("a." + testSerial).Message("hello_world").UsePost(true).Execute()
	}
	runtime.Gosched()
	//time.Sleep(time.Millisecond * 10)
	cancelPub()
	// pn.UnsubscribeAll()
	// pn.RemoveListener(listener)
	// pn.Destroy()

	//fmt.Println("after cancelPub")
	pn.UnsubscribeAll()
	//fmt.Println("after UnsubscribeAll")
	pn.RemoveListener(listener)
	//fmt.Println("after RemoveListener")
	pn.Destroy()
	//fmt.Println("after Destroy")

	pn = nil
}

func TestDestroy2(t *testing.T) {
	testSerial := "bb"
	config := configCopy()
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	if enableDebuggingInTests {
		config.Log = log.New(os.Stdout, "", log.Ltime|log.Lmicroseconds)
	}
	config.WPSReconnectionPolicy = webpubsub.WPSExponentialPolicycy
	config.MaximumReconnectionRetries = -1
	config.UUID = testSerial
	config.SuppressLeaveEvents = true
	config.SetPresenceTimeoutWithCustomInterval(330, 300)
	config.MaxWorkers = 0

	pn := webpubsub.NewWebPubSub(config)
	listener := webpubsub.NewListener()
	pn.AddListener(listener)
	pn.Subscribe().Channels([]string{"a." + testSerial}).Execute()
	<-listener.Status
	pn.Destroy()
}
