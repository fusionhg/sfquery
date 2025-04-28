# sfquery

Allows a Go app to act as a Connected App and run SOQL queries against a Salesforce instance API.  Can automatically loop through multipage result sets from Salesforce and collate into a single slice of returned records.

## Salesforce Setup

* Create a new Connected App
* Enable OAuth Settings
* Enable for Device Flow (may or may not be required)
* Note the produced "Consumer Key" (Client ID) and "Consumer Secret" (Client Secret)

## Basic Use

Install—
```bash
GOPROXY=direct go get github.com/fusionhg/sfquery@latest
```

Acquire access token—
```go
rtOut, err = sfquery.RequestOauth2Token(
    "{yourinstance}.salesforce.com",
    sfquery.RequestOauth2TokenIn{
        GrantType:    "client_credentials",
        ClientId:     "{your Client ID here}",
        ClientSecret: "{your Client Secret here}",
    },
)
if err != nil {
    panic(err)
}
accessToken := rtOut.AccessToken
```

Run query—
```go
type TestContact struct {
	Id       string
	LastName string
}

records, err := sfquery.TypedQueryMultipage[TestContact](
    "{yourinstance}.salesforce.com",
    accessToken,
    "select Id, LastName from Contact",
)
if err != nil {
    t.Fatal(err)
}
```
