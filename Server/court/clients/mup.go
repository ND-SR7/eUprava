package clients

import (
	"bytes"
	"context"
	"court/data"
	"court/domain"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MUPClient struct {
	client  *http.Client
	address string
}

func NewMUPClient(client *http.Client, address string) MUPClient {
	return MUPClient{
		client:  client,
		address: address,
	}
}

// Client methods

// Notifies MUP to create driving ban based on created suspension
func (mc *MUPClient) NotifyOfSuspension(ctx context.Context, newSuspension data.NewSuspension, token string) error {
	toDateTime, err := time.Parse("2006-01-02T15:04:05", newSuspension.To)
	if err != nil {
		return err
	}

	personID, err := primitive.ObjectIDFromHex(newSuspension.Person)
	if err != nil {
		return err
	}

	drivingBan := data.DrivingBan{
		Reason:   "License suspension",
		Duration: toDateTime,
		Person:   personID,
	}

	requestBody, err := json.Marshal(drivingBan)
	if err != nil {
		return err
	}

	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mc.address+"/driving-ban", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := mc.client.Do(req)
	if err != nil {
		return handleHttpReqErr(err, mc.address+"/driving-ban", http.MethodPost, timeout)
	}

	if resp.StatusCode != http.StatusCreated {
		return domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}
