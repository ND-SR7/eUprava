package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"statistics/data"
)

type MupClient struct {
	client  *http.Client
	address string
}

func NewMupClient(client *http.Client, address string) MupClient {
	return MupClient{
		client:  client,
		address: address,
	}
}

func (c *MupClient) GetAllRegisteredVehicles(ctx context.Context) ([]data.Vehicle, error) {
	url := c.address + "/api/v1/registered-vehicles"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var vehicles []data.Vehicle
	if err := json.NewDecoder(resp.Body).Decode(&vehicles); err != nil {
		return nil, err
	}

	return vehicles, nil
}
