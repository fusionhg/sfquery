package sfquery

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type TestContact struct {
	Id       string
	LastName string
}

func TestQuery(t *testing.T) {
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
		t.Fatal(err)
	}
	qres, err := TypedQuery[TestContact](
		os.Getenv("SF_DOMAIN_NAME"),
		rtOut.AccessToken,
		"select Id, LastName from Contact",
	)
	if err != nil {
		t.Fatal(err)
	}
	assert.Greater(t, len(qres.Records), 0)
	assert.NotEmpty(t, qres.Records[0].Id)
	assert.Equal(t, qres.Records[0].LastName, "Benioff (Sample)")
}

func TestQueryMultipage(t *testing.T) {
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
		t.Fatal(err)
	}
	qres, err := TypedQueryMultipage[TestContact](
		os.Getenv("SF_DOMAIN_NAME"),
		rtOut.AccessToken,
		"select Id, LastName from Contact",
		4000,
	)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, qres)
	records := *qres
	// this will only pass if you have >= 4000 contacts in salesforce
	assert.Equal(t, len(records), 4000)
	assert.NotEmpty(t, records[0].Id)
	assert.NotEqual(t, records[0].Id, records[2000].Id)
}
