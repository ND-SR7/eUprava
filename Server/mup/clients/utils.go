package clients

import (
	"fmt"
	"math/rand"
	"mup/domain"
	"net/url"
	"time"
)

func handleHttpReqErr(err error, reqUrl string, method string, timeout time.Duration) error {
	urlErr, ok := err.(*url.Error)
	if !ok {
		return domain.ErrUnknown{
			InnerErr: err,
		}
	}
	if urlErr.Timeout() {
		return domain.ErrClientSideTimeout{
			URL:        reqUrl,
			Method:     method,
			MaxTimeout: timeout,
		}
	}
	return domain.ErrConnecting{
		Err: urlErr,
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GeneratePlates() string {
	return fmt.Sprintf("SRB-%s-%s", RandString(3), RandString(2))
}

func GenerateRegistration() string {
	return RandString(8)
}
