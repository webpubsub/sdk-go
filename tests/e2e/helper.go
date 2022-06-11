package e2e

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/stretchr/testify/assert"
	webpubsub "github.com/webpubsub/go/v7"
)

var enableDebuggingInTests = false

const (
	SPECIAL_CHARACTERS = "-.,_~:/?#[]@!$&'()*+;=`|"
	SPECIAL_CHANNEL    = "-._~:/?#[]@!$&'()*+;=`|"
)

var pamConfig *webpubsub.Config
var config *webpubsub.Config

var (
	serverErrorTemplate     = "webpubsub/server: Server respond with error code %d"
	validationErrorTemplate = "webpubsub/validation: %s"
	connectionErrorTemplate = "webpubsub/connection: %s"
)

func seedRand() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func init() {
	seedRand()
	config = webpubsub.NewConfig(webpubsub.GenerateUUID())
	config.PublishKey = os.Getenv("PUBLISH_KEY")
	config.SubscribeKey = os.Getenv("SUBSCRIBE_KEY")

	pamConfig = webpubsub.NewConfig(webpubsub.GenerateUUID())
	pamConfig.PublishKey = os.Getenv("PAM_PUBLISH_KEY")
	pamConfig.SubscribeKey = os.Getenv("PAM_SUBSCRIBE_KEY")
	pamConfig.SecretKey = os.Getenv("PAM_SECRET_KEY")
}

func configCopy() *webpubsub.Config {
	cfg := new(webpubsub.Config)
	*cfg = *config
	cfg.UUID = webpubsub.GenerateUUID()
	return cfg
}

func pamConfigCopy() *webpubsub.Config {
	config := new(webpubsub.Config)
	*config = *pamConfig
	config.UUID = webpubsub.GenerateUUID()
	return config
}

func randomized(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, rand.Intn(10000000))
}

type fakeTransport struct {
	Status     string
	StatusCode int
	Body       io.ReadCloser
}

func (t fakeTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     t.Status,
		StatusCode: t.StatusCode,
		Body:       t.Body,
	}, nil
}

func (t fakeTransport) Dial(string, string) (net.Conn, error) {
	return nil, errors.New("ooops!")
}

func logInTest(format string, a ...interface{}) (n int, err error) {
	if enableDebuggingInTests {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

func checkFor(assert *assert.Assertions, maxTime, intervalTime time.Duration, fun func() error) {
	maxTimeout := time.NewTimer(maxTime)
	interval := time.NewTicker(intervalTime)
	lastErr := fun()
	if lastErr == nil {
		return
	}
ForLoop:
	for {
		select {
		case <-interval.C:
			lastErr := fun()
			if lastErr != nil {
				logInTest("Error: %s. Checking in next %s\n", lastErr, intervalTime)
				continue
			} else {
				break ForLoop
			}
		case <-maxTimeout.C:
			assert.Fail(lastErr.Error())
			break ForLoop
		}
	}
}

func heyIterator(count int) <-chan string {
	channel := make(chan string)

	init := "hey-"

	go func() {
		for i := 1; i <= count; i++ {
			channel <- fmt.Sprintf("%s%d", init, i)
		}
	}()

	return channel
}
