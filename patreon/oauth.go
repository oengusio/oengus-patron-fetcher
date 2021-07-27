package patreon

import (
    "bytes"
    "encoding/json"
    "errors"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "oenugs-patreon/structs"
    "os"
    "strings"
)

func Oauth2FetchToken(code string) (structs.PatreonTokens, error) {
    var response structs.PatreonTokens

    query := url.Values{}
    query.Set("grant_type", "authorization_code")
    query.Set("code", code)
    query.Set("client_id", os.Getenv("PATREON_CLIENT_ID"))
    query.Set("client_secret", os.Getenv("PATREON_CLIENT_SECRET"))
    query.Set("redirect_uri", os.Getenv("OENGUS_BASE") + "/user/settings/sync/patreon")

    queryBytes := []byte(query.Encode())

    apiUrl := "https://api.patreon.com/oauth2/token"
    req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(queryBytes))

    if err != nil {
        log.Println(err)
        return response, err
    }

    req.Header.Set("User-Agent", "oengus.io/patreon-fetcher")

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

    return response, nil
}

func Oauth2FetchUser(token structs.PatreonTokens) {
    apiUrl := "https://api.patreon.com/oauth2/v2/identity"
    req, err := http.NewRequest(http.MethodPost, apiUrl, nil)

    if err != nil {
        log.Println(err)
        return response, err
    }

    req.Header.Set("User-Agent", "oengus.io/patreon-fetcher")
}