package sfquery

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
)

func callUrlWithToken(u string, bearerToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received status code %d: %s", resp.StatusCode, string(body))
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

type TypedQueryResult[T any] struct {
	TotalSize      int    `json:"totalSize"`
	Done           bool   `json:"done"`
	NextRecordsURL string `json:"nextRecordsUrl"`
	Records        []T    `json:"records"`
}

// pull just the first page of results back from salesforce
// (max 2000 results)
func TypedQuery[T any](sfDomainName string, accessToken string, queryOrNextRecordsPath string) (*TypedQueryResult[T], error) {
	u := url.URL{
		Scheme: "https",
		Host:   sfDomainName,
	}
	if strings.HasPrefix(queryOrNextRecordsPath, "/") {
		u.Path = queryOrNextRecordsPath
	} else {
		u.Path = fmt.Sprintf("/services/data/v%s/query", SfApiVersion)
		v := url.Values{}
		v.Set("q", queryOrNextRecordsPath)
		u.RawQuery = v.Encode()
	}
	body, err := callUrlWithToken(u.String(), accessToken)
	if err != nil {
		return nil, err
	}
	var out TypedQueryResult[T]
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// loop to pull back multiple pages of results
func TypedQueryMultipage[T any](sfDomainName string, accessToken string, query string, maxRecords int) (*[]T, error) {
	var out []T
	queryOrNextRecordsPath := query
	var qres *TypedQueryResult[T]
	var err error
	for (qres == nil || !qres.Done) && (len(out) < maxRecords) {
		qres, err = TypedQuery[T](sfDomainName, accessToken, queryOrNextRecordsPath)
		if err != nil {
			return nil, err
		}
		out = append(out, qres.Records...)
		queryOrNextRecordsPath = qres.NextRecordsURL
	}
	return &out, nil
}

func TypedQueryAll[T any](sfDomainName string, accessToken string, query string) (*[]T, error) {
	return TypedQueryMultipage[T](sfDomainName, accessToken, query, math.MaxInt)
}
