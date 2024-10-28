package paymentgatewaya

import (
	"bytes"
	"encoding/json"
	"net/http"
	"seta/pkg/clients/paymentgateway"
	"seta/pkg/model"
	"time"

	"github.com/shopspring/decimal"
)

type Client struct {
	Endpoint   string
	HTTPClient *http.Client
	// AuthToken string // implied that the client requires an auth token
}

func ClientProvider(Endpoint string) paymentgateway.IPaymentGateway {
	httpClient := &http.Client{}
	httpClient.Timeout = 60 * time.Second

	return &Client{
		Endpoint:   Endpoint,
		HTTPClient: httpClient,
	}
}

func (c *Client) Deposit(AccountID string, amount decimal.Decimal) (*model.TransactionResponse, *int, error) {
	depositRequest := &model.DepositRequest{
		AccountID: AccountID,
		Amount:    amount,
	}

	payload, err := json.Marshal(depositRequest)
	if err != nil {
		return nil, nil, err
	}

	// Call the payment gateway API
	resp, err := c.HTTPClient.Post(c.Endpoint+"/deposit", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, &resp.StatusCode, err
	}

	// Parse the response
	var gatewayATransactionResponse model.GatewayATransactionResponse
	err = json.NewDecoder(resp.Body).Decode(&gatewayATransactionResponse)
	if err != nil {
		return nil, nil, err
	}

	transactionResponse := model.MapGatewayATransactionResponse(&gatewayATransactionResponse)

	return &transactionResponse, &resp.StatusCode, nil
}

func (c *Client) Withdraw(AccountID string, amount decimal.Decimal) (*model.TransactionResponse, *int, error) {
	withdrawRequest := &model.WithdrawRequest{
		AccountID: AccountID,
		Amount:    amount,
	}

	payload, err := json.Marshal(withdrawRequest)
	if err != nil {
		return nil, nil, err
	}

	// Call the payment gateway API
	resp, err := c.HTTPClient.Post(c.Endpoint+"/withdraw", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, &resp.StatusCode, err
	}

	// Parse the response
	var gatewayATransactionResponse model.GatewayATransactionResponse
	err = json.NewDecoder(resp.Body).Decode(&gatewayATransactionResponse)
	if err != nil {
		return nil, nil, err
	}

	transactionResponse := model.MapGatewayATransactionResponse(&gatewayATransactionResponse)

	return &transactionResponse, &resp.StatusCode, nil
}
