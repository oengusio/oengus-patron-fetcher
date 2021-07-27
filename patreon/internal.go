package patreon

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "oenugs-patreon/structs"
    "os"
    "strings"
)

var campaignId = os.Getenv("PATREON_CAMPAIGN_ID")
var httpClient = http.Client{}

func RefreshToken(tokens structs.PatreonTokens) (structs.PatreonTokens, error) {
    var response structs.PatreonTokens

    apiUrl := "https://api.patreon.com/oauth2/token"
    req, err := http.NewRequest(http.MethodPost, apiUrl, nil)
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

    req.URL.RawQuery = query.Encode()

    res, httpErr := httpClient.Do(req)
    if httpErr != nil {
        log.Println(httpErr)
        return response, httpErr
    }

    if res.Body != nil {
        defer res.Body.Close()
    }

    body, readErr := ioutil.ReadAll(res.Body)
    if readErr != nil {
        log.Println(readErr)
        return response, readErr
    }

    // In case I ever need the new tokens
    strBody := string(body)
    log.Println("credentials fetched: " + strBody)

    if strings.Contains(strBody, "error") {
        return response, errors.New(strBody)
    }

    jsonErr := json.Unmarshal(body, &response)
    if jsonErr != nil {
        log.Println(jsonErr)
        return response, jsonErr
    }

    file, _ := json.MarshalIndent(response, "", " ")
    _ = ioutil.WriteFile("/storage/oengus-patreon/patreon-credentials.json", file, 0644)

    return response, nil
}

func FetchPatrons(tokens structs.PatreonTokens) (structs.PatreonMembersResponse, error) {
    var response structs.PatreonMembersResponse

    apiUrl := fmt.Sprintf("https://api.patreon.com/oauth2/v2/campaigns/%s/members", campaignId)
    req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
    if err != nil {
        log.Println(err)
        return response, err
    }

    req.Header.Set("User-Agent", "oengus.io/patreon-fetcher")
    req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

    query := req.URL.Query()
    query.Set("fields[member]", "full_name,patron_status,will_pay_amount_cents")
    query.Set("include", "user")
    query.Set("page[count]", "1000")

    req.URL.RawQuery = query.Encode()

    res, httpErr := httpClient.Do(req)
    if httpErr != nil {
        log.Println(httpErr)
        return response, httpErr
    }

    // defer calls are not executed until the function returns
    if res.Body != nil {
        defer res.Body.Close()
    }

    if res.StatusCode == http.StatusUnauthorized {
        return response, errors.New("StatusUnauthorized")
    }

    body, readErr := ioutil.ReadAll(res.Body)
    if readErr != nil {
        log.Print(readErr)
        return response, readErr
    }

    log.Println(string(body))

    jsonErr := json.Unmarshal(body, &response)
    if jsonErr != nil {
        log.Println(jsonErr)
        return response, jsonErr
    }

    return response, nil
}
