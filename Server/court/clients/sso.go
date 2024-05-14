package clients

import (
	"context"
	"court/data"
	"court/domain"
	"encoding/json"
	"fmt"
	"net/http"
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

// Retrieves person based on provided account ID
func (sc *SSOClient) GetPersonByID(ctx context.Context, accountID, token string) (data.Person, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sc.address+"/user/"+accountID, nil)
	if err != nil {
		return data.Person{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := sc.client.Do(req)
	if err != nil {
		return data.Person{}, handleHttpReqErr(err, sc.address+"/user/"+accountID, http.MethodGet, timeout)
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

// Retrieves person based on provided email
func (sc *SSOClient) GetPersonByEmail(ctx context.Context, email, token string) (data.Person, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sc.address+"/user/email/"+email, nil)
	if err != nil {
		return data.Person{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := sc.client.Do(req)
	if err != nil {
		return data.Person{}, handleHttpReqErr(err, sc.address+"/user/email/"+email, http.MethodGet, timeout)
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

// Retrieves legal entity based on provided account ID
func (sc *SSOClient) GetLegalEntityByID(ctx context.Context, accountID, token string) (data.LegalEntity, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sc.address+"/user/"+accountID, nil)
	if err != nil {
		return data.LegalEntity{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := sc.client.Do(req)
	if err != nil {
		return data.LegalEntity{}, handleHttpReqErr(err, sc.address+"/user/"+accountID, http.MethodGet, timeout)
	}

	if resp.StatusCode != http.StatusOK {
		return data.LegalEntity{}, domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	var legalEntity data.LegalEntity
	if err := json.NewDecoder(resp.Body).Decode(&legalEntity); err != nil {
		return data.LegalEntity{}, fmt.Errorf("failed to decode JSON response: %s", err.Error())
	}

	return legalEntity, nil
}

// Retrieves legal entity based on provided email
func (sc *SSOClient) GetLegalEntityByEmail(ctx context.Context, email, token string) (data.LegalEntity, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sc.address+"/user/email/"+email, nil)
	if err != nil {
		return data.LegalEntity{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := sc.client.Do(req)
	if err != nil {
		return data.LegalEntity{}, handleHttpReqErr(err, sc.address+"/user/email/"+email, http.MethodGet, timeout)
	}

	if resp.StatusCode != http.StatusOK {
		return data.LegalEntity{}, domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	var legalEntity data.LegalEntity
	if err := json.NewDecoder(resp.Body).Decode(&legalEntity); err != nil {
		return data.LegalEntity{}, fmt.Errorf("failed to decode JSON response: %s", err.Error())
	}

	return legalEntity, nil
}
