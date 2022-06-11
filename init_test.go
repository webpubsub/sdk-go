package webpubsub

var pnconfig *Config
var webpubsub *WebPubSub

func init() {
	pnconfig = NewConfig(GenerateUUID())

	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"

	webpubsub = NewWebPubSub(pnconfig)
}

func webpubsubCopy() *WebPubSub {
	pn := new(WebPubSub)
	*pn = *webpubsub
	return pn
}
