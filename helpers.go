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

func RefreshToken() {
	//
}

// TODO: https://www.reddit.com/r/golang/comments/2xmnvs/returning_nil_for_a_struct/
func FetchPatrons() *PatreonMembersResponse {
	url := fmt.Sprintf("https://www.patreon.com/api/oauth2/v2/campaigns/%s/members", campaignId)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)

		return nil
	}

	req.Header.Set("User-Agent", "oengus.io/patreon-fetcher")

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	parsedBody := &PatreonMembersResponse{}

	err2 := json.Unmarshal(body, &parsedBody)
	if err2 != nil {
		fmt.Println(err2)
		return nil
	}

	return parsedBody
}