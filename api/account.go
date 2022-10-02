package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// AccountInformation model
type AccountInformation struct {
	// Account credit balance display text
	AccountBalance string `json:"accountBalance,omitempty"`

	// Account credit balance in cents
	AccountBalanceRaw string `json:"accountBalanceRaw,omitempty"`

	// Account credit balance currency
	AccountBalanceRawCurrency string `json:"accountBalanceRawCurrency,omitempty"`

	// Company name
	CompanyName string `json:"companyName,omitempty"`

	// User is part of a team
	IsTeamMember bool `json:"isTeamMember,omitempty"`

	// Team leader email
	TeamLeader string `json:"teamLeader,omitempty"`

	// User Avatar URL
	UserAvatar string `json:"userAvatar,omitempty"`

	// User email
	UserEmail string `json:"userEmail,omitempty"`

	// User ID
	UserId int64 `json:"userId,omitempty"`

	// User name
	UserName string `json:"userName,omitempty"`
}

func (c *Client) GetAccountInformation(ctx context.Context) (*AccountInformation, error) {
	req, err := newGetAccountInformationRequest(c.Server)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if errorStatus(res.StatusCode) {
		return nil, fmt.Errorf("error getting account information: %s", res.Status)
	}

	decoder := json.NewDecoder(res.Body)

	defer res.Body.Close()

	account := &AccountInformation{}

	if err := decoder.Decode(account); err != nil {
		return nil, err
	}

	return account, nil
}

// newGetAccountInformationRequest generates requests for GetAccountInformation
func newGetAccountInformationRequest(server string) (*http.Request, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "account/accountInformation"

	return http.NewRequest("GET", serverURL.String(), nil)
}
