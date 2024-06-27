package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"police/data"

	"go.mongodb.org/mongo-driver/mongo"
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

func (mc MupClient) GetRegistrationByPlate(ctx context.Context, plates data.PlateRequest, token string) (data.Registration, error) {
	requestBody, err := json.Marshal(plates)
	if err != nil {
		return data.Registration{}, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mc.address+"/registration-by-plate", bytes.NewBuffer(requestBody))
	if err != nil {
		return data.Registration{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := mc.client.Do(req)
	if err != nil {
		return data.Registration{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return data.Registration{}, errors.New("unexpected status code: " + resp.Status)
	}

	var registration data.Registration
	if err := json.NewDecoder(resp.Body).Decode(&registration); err != nil {
		return data.Registration{}, err
	}

	return registration, nil
}

func (mc MupClient) CheckDrivingBan(ctx context.Context, jmbg data.JMBGRequest, token string) (*data.DrivingBan, error) {
	requestBody, err := json.Marshal(jmbg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mc.address+"/check-persons-driving-ban", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := mc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	var drivingBan data.DrivingBan
	err = json.NewDecoder(resp.Body).Decode(&drivingBan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &drivingBan, nil
}

func (mc MupClient) GetDrivingPermitByJMBG(ctx context.Context, jmbg data.JMBGRequest, token string) (data.TrafficPermit, error) {
	var permit data.TrafficPermit

	requestBody, err := json.Marshal(jmbg)
	if err != nil {
		return permit, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mc.address+"/check-persons-driving-permit", bytes.NewBuffer(requestBody))
	if err != nil {
		return permit, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := mc.client.Do(req)
	if err != nil {
		return permit, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return permit, errors.New("unexpected status code: " + resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&permit); err != nil {
		return permit, err
	}

	return permit, nil
}
