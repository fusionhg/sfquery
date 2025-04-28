package sfquery

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestClientCredentialsAuthBad(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("please copy .env.example to .env and fill with connected app credentials"))
	}
	_, err = RequestOauth2Token(
		os.Getenv("SF_DOMAIN_NAME"),
		RequestOauth2TokenIn{
			GrantType:    "client_credentials",
			ClientId:     "notAClientId",
			ClientSecret: "notAClientSecret",
		},
	)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "received status code 40")
}

func TestClientCredentialsAuthOk(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("please copy .env.example to .env and fill with connected app credentials"))
	}
	rtOut, err := RequestOauth2Token(
		os.Getenv("SF_DOMAIN_NAME"),
		RequestOauth2TokenIn{
			GrantType:    "client_credentials",
			ClientId:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		},
	)
	if err != nil {
		panic(err)
	}
	assert.NotEmpty(t, rtOut.AccessToken)
	// log.Printf("access token: %s", rtOut.AccessToken)
}
