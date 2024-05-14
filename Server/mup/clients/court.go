package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mup/data"
	"mup/domain"
	"net/http"
	"time"
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

func (cc CourtClient) CheckForPersonsWarrant(ctx context.Context, userID primitive.ObjectID, token string) (data.Warrants, error) {
	requestBody, err := json.Marshal(userID.Hex())
	if err != nil {
		_ = fmt.Errorf("failed to marshal user id: %v", err)
		return data.Warrants{}, err
	}

	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cc.address+"/warrants/"+userID.Hex(), bytes.NewBuffer(requestBody))
	if err != nil {
		return data.Warrants{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := cc.client.Do(req)
	if err != nil {
		return data.Warrants{}, handleHttpReqErr(err, cc.address, http.MethodPost, timeout)
	}

	if resp.StatusCode != http.StatusOK {
		return data.Warrants{}, domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	var serviceResponse data.Warrants
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&serviceResponse); err != nil {
		return data.Warrants{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return serviceResponse, nil
}
