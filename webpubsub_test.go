package webpubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializer(t *testing.T) {
	assert := assert.New(t)

	pnconfig := NewConfig(GenerateUUID())
	pnconfig.PublishKey = "my_pub_key"
	pnconfig.SubscribeKey = "my_sub_key"
	pnconfig.SecretKey = "my_secret_key"
	webpubsub := NewWebPubSub(pnconfig)

	assert.Equal("my_pub_key", webpubsub.Config.PublishKey)
	assert.Equal("my_sub_key", webpubsub.Config.SubscribeKey)
	assert.Equal("my_secret_key", webpubsub.Config.SecretKey)
}

func TestDemoInitializer(t *testing.T) {
	demo := NewWebPubSubDemo()

	assert := assert.New(t)

	assert.Equal("demo", demo.Config.PublishKey)
	assert.Equal("demo", demo.Config.SubscribeKey)
	assert.Equal("demo", demo.Config.SecretKey)
}

func TestMultipleConcurrentInit(t *testing.T) {
	c1 := NewConfig(GenerateUUID())
	go NewWebPubSub(c1)
	c2 := NewConfig(GenerateUUID())
	NewWebPubSub(c2)
}
