package sfquery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type RequestOauth2TokenIn struct {
	GrantType    string
	ClientId     string
	ClientSecret string
}

func (req *RequestOauth2TokenIn) ToValues() url.Values {
	v := url.Values{}
	v.Set("grant_type", req.GrantType)
	v.Set("client_id", req.ClientId)
	v.Set("client_secret", req.ClientSecret)
	return v
}

type RequestOauth2TokenOut struct {
	AccessToken string `json:"access_token"`
	Signature   string `json:"signature"`
	InstanceUrl string `json:"instance_url"`
	Id          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
}

// https://help.salesforce.com/s/articleView?id=xcloud.remoteaccess_oauth_client_credentials_flow.htm&type=5
func RequestOauth2Token(sfDomainName string, in RequestOauth2TokenIn) (*RequestOauth2TokenOut, error) {
	reqBody := in.ToValues().Encode()
	url := url.URL{
		Scheme: "https",
		Host:   sfDomainName,
		Path:   "/services/oauth2/token",
	}
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var out RequestOauth2TokenOut
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
