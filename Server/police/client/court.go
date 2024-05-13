package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"police/data"
)

type CourtClient struct {
	client  *http.Client
	address string
}

func NewCourtClient(client *http.Client, address string) CourtClient {
	return CourtClient{
		client:  client,
		address: address,
	}
}

func (cc CourtClient) CreateCrimeReport(ctx context.Context, violation data.TrafficViolation) (interface{}, error) {
	requestBody, err := json.Marshal(violation)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cc.address, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := cc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	var response interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}
