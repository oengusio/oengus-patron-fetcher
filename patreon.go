package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var campaignId = os.Getenv("PATREON_CAMPAIGN_ID")
var httpClient = http.Client{}


/// <editor-fold name="Refresh token response">
type PatreonTokens struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
/// </editor-fold>

/// <editor-fold name="Patrons response">
type PatreonMembersResponse struct {
	Data []PatreonMembersData `json:"data"`
}

type PatreonMembersData struct {
	Attributes PatreonMembersAttribute `json:"attributes"`
}

type PatreonMembersAttribute struct {
	FullName string `json:"full_name"`
	PatronStatus string `json:"patron_status"`
	WillPayAmountCents int `json:"will_pay_amount_cents"`
}
/// </editor-fold>

func RefreshToken(tokens PatreonTokens) (PatreonTokens, error) {
	// TODO: test if initial refresh token stays the same
	var response PatreonTokens

	url := "https://www.patreon.com/api/oauth2/token"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
		return response, err
	}

	req.Header.Set("User-Agent", "oengus.io/patreon-fetcher")

	query := req.URL.Query()
	query.Set("grant_type", "refresh_token")
	query.Set("refresh_token", tokens.RefreshToken)
	query.Set("client_id", os.Getenv("PATREON_CLIENT_ID"))
	query.Set("client_secret", os.Getenv("PATREON_CLIENT_SECRET"))

	res, httpErr := httpClient.Do(req)
	if httpErr != nil {
		log.Fatal(httpErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		return response, jsonErr
	}

	return response, nil
}

// TODO: https://www.reddit.com/r/golang/comments/2xmnvs/returning_nil_for_a_struct/
func FetchPatrons(tokens PatreonTokens) (PatreonMembersResponse, error) {
	var response PatreonMembersResponse

	url := fmt.Sprintf("https://www.patreon.com/api/oauth2/v2/campaigns/%s/members", campaignId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
		return response, err
	}

	req.Header.Set("User-Agent", "oengus.io/patreon-fetcher")
	req.Header.Set("Authorisation", "Bearer " + tokens.AccessToken)

	query := req.URL.Query()
	query.Set("fields%5Bmember%5D", "full_name,patron_status,will_pay_amount_cents")
	query.Set("&page%5Bcount%5D", "1000")

	res, httpErr := httpClient.Do(req)
	if httpErr != nil {
		log.Fatal(httpErr)
	}

	// defer calls are not executed until the function returns
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		return response, jsonErr
	}

	return response, nil
}