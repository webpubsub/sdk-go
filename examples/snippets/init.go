package main

var pn *WebPubSub

func Init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "demo"
	pnconfig.SubscribeKey = "demo"

	pn = NewWebPubSub(pnconfig)
}

func webpubsubCopy() *WebPubSub {
	pn := new(WebPubSub)
	*pn = *webpubsub
	return pn
}
