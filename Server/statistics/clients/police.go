package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"statistics/data"
)

type PoliceClient struct {
	client  *http.Client
	address string
}

func NewPoliceClient(client *http.Client, address string) PoliceClient {
	return PoliceClient{
		client:  client,
		address: address,
	}
}

func (pc *PoliceClient) GetTrafficViolations(ctx context.Context, token string) (data.TrafficViolations, error) {
	url := pc.address + "/traffic-violation"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := pc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve traffic violations: %s", resp.Status)
	}

	var violations data.TrafficViolations
	if err := json.NewDecoder(resp.Body).Decode(&violations); err != nil {
		return nil, err
	}

	return violations, nil
}
