package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"police/data"
	"police/domain"
	"time"
)

type SSOClient struct {
	client  *http.Client
	address string
}

func NewSSOClient(client *http.Client, address string) SSOClient {
	return SSOClient{
		client:  client,
		address: address,
	}
}

// Client methods

// Retrieves person based on provided JMBG
func (sc *SSOClient) GetPersonByJMBG(ctx context.Context, jmbg, token string) (data.Person, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sc.address+"/user/jmbg/"+jmbg, nil)
	if err != nil {
		return data.Person{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := sc.client.Do(req)
	if err != nil {
		return data.Person{}, handleHttpReqErr(err, sc.address+"/user/jmbg/"+jmbg, http.MethodGet, timeout)
	}

	if resp.StatusCode != http.StatusOK {
		return data.Person{}, domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	var person data.Person
	if err := json.NewDecoder(resp.Body).Decode(&person); err != nil {
		return data.Person{}, fmt.Errorf("failed to decode JSON response: %s", err.Error())
	}

	return person, nil
}
