package clients

import (
	"net/url"
	"police/domain"
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
