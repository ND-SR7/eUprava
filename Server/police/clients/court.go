package clients

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

func (cc CourtClient) CreateCrimeReport(ctx context.Context, violation data.TrafficViolation, token string) error {
	requestBody, err := json.Marshal(violation)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cc.address+"/crime-report", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := cc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + resp.Status)
	}

	return nil
}
