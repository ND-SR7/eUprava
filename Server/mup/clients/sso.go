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

//Client methods (checking username and ID)

func (ssoc SSOClient) GetUserIdByEmail(ctx context.Context, email, token string) (string, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ssoc.address+"/users/email/"+email, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := ssoc.client.Do(req)
	if err != nil {
		return "", handleHttpReqErr(err, ssoc.address+"/users/email/"+email, http.MethodPost, timeout)
	}

	if resp.StatusCode != http.StatusOK {
		return "", domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	// Parse the JSON response
	var serviceResponse data.Person
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&serviceResponse); err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return serviceResponse.Account.ID.Hex(), nil
}

func (ssoc SSOClient) GetUserById(ctx context.Context, id primitive.ObjectID, token string) (data.Person, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	idUser := data.UserId{}
	idUser.ID = id
	requestBody, err := json.Marshal(idUser)
	if err != nil {
		return data.Person{}, fmt.Errorf("failed to marshal id: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ssoc.address+"/get-user-by-id", bytes.NewBuffer(requestBody))
	if err != nil {
		return data.Person{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := ssoc.client.Do(req)
	if err != nil {
		return data.Person{}, handleHttpReqErr(err, ssoc.address+"/get-user-by-id", http.MethodPost, timeout)
	}

	if resp.StatusCode != http.StatusOK {
		return data.Person{}, domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	// Parse the JSON response
	var serviceResponse data.Person
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&serviceResponse); err != nil {
		return data.Person{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return serviceResponse, nil
}
