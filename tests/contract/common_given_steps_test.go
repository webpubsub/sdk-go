package contract

import (
	"context"

	webpubsub "github.com/webpubsub/sdk-go/v7"
)

func iHaveAKeysetWithAccessManagerEnabled(ctx context.Context) error {
	state := getCommonState(ctx)
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.PublishKey = state.contractTestConfig.publishKey
	config.SubscribeKey = state.contractTestConfig.subscribeKey
	config.SecretKey = state.contractTestConfig.secretKey
	config.Origin = state.contractTestConfig.hostPort
	config.Secure = state.contractTestConfig.secure

	state.WebPubSub = webpubsub.NewWebPubSub(config)
	return nil
}

func iHaveAKeysetWithAccessManagerEnabledWithoutSecretKey(ctx context.Context) error {
	state := getCommonState(ctx)
	config := webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.PublishKey = state.contractTestConfig.publishKey
	config.SubscribeKey = state.contractTestConfig.subscribeKey
	config.Origin = state.contractTestConfig.hostPort
	config.Secure = state.contractTestConfig.secure

	state.WebPubSub = webpubsub.NewWebPubSub(config)
	return nil
}
