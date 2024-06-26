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

func (mc MupClient) GetRegistrationByPlate(ctx context.Context, plates data.Plates, token string) error {
	requestBody, err := json.Marshal(plates)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mc.address+"/registration-by-plate", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := mc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + resp.Status)
	}

	return nil
}

func (mc MupClient) CheckDrivigBan(ctx context.Context, jmbg data.JMBGRequest, token string) (bool, error) {
	requestBody, err := json.Marshal(jmbg)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mc.address+"/check-persons-driving-ban", bytes.NewBuffer(requestBody))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := mc.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("unexpected status code: " + resp.Status)
	}

	var result struct {
		DrivingBan bool `json:"drivingBan"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return result.DrivingBan, nil
}

func (mc MupClient) GetDrivingPermitByJMBG(ctx context.Context, jmbg data.JMBGRequest, token string) (data.TrafficPermit, error) {
	var permit data.TrafficPermit

	requestBody, err := json.Marshal(jmbg)
	if err != nil {
		return permit, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mc.address+"/check-for-persons-driving-permit", bytes.NewBuffer(requestBody))
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

func (mc MupClient) GetVehicleRegistration(ctx context.Context, CheckVehicleRegistration data.CheckVehicleRegistration, token string) (data.Registration, error) {
	var registration data.Registration

	platesNumber := data.PlateRequest{
		Plate: CheckVehicleRegistration.PlatesNumber,
	}

	requestBody, err := json.Marshal(platesNumber)
	if err != nil {
		return registration, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mc.address+"/registration-by-plate", bytes.NewBuffer(requestBody))
	if err != nil {
		return registration, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := mc.client.Do(req)
	if err != nil {
		return registration, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return registration, errors.New("unexpected status code: " + resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&registration); err != nil {
		return registration, err
	}

	return registration, nil
}
