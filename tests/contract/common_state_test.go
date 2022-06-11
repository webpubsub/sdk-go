package contract

import webpubsub "github.com/webpubsub/go/v7"

type commonStateKey struct{}

type commonState struct {
	contractTestConfig contractTestConfig
	WebPubSub          *webpubsub.WebPubSub
	err                error
	statusResponse     webpubsub.StatusResponse
}

func newCommonState(contractTestConfig contractTestConfig) *commonState {

	return &commonState{
		contractTestConfig: contractTestConfig,
	}
}
