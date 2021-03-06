package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	webpubsub "github.com/webpubsub/sdk-go/v7"
	"github.com/webpubsub/sdk-go/v7/tests/stubs"
)

func TestTime(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/time/0",
		Query:              "",
		ResponseBody:       `[15078947309567840]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	pn := webpubsub.NewWebPubSub(configCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.Time().Execute()

	assert.Nil(err)

	assert.True(int64(15059085932399340) < res.Timetoken)
}

func TestTimeContext(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/time/0",
		Query:              "",
		ResponseBody:       `[15078947309567840]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	pn := webpubsub.NewWebPubSub(configCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.TimeWithContext(backgroundContext).Execute()

	assert.Nil(err)

	assert.True(int64(15059085932399340) < res.Timetoken)
}
