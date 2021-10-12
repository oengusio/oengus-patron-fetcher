package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "oenugs-patreon/cache"
    "oenugs-patreon/patreon"
    "oenugs-patreon/structs"
    "os"
    "os/signal"
    "time"
)

func UpdatePatrons() {
    log.Println("Updating patrons")

    patrons, err := patreon.FetchPatrons(cache.Tokens)

    // 401 response, refresh the tokens
    if err != nil && err.Error() == "StatusUnauthorized" {
        newTokens, refreshErr := patreon.RefreshToken(cache.Tokens)

        if refreshErr != nil {
            log.Println("Error while refreshing tokens", refreshErr.Error())
            return
        }

        cache.Tokens = newTokens
        patrons, err = patreon.FetchPatrons(cache.Tokens)
    }

    if err != nil {
        log.Println(err)
        return
    }

    newCache := structs.PatronOutput{
        Patrons: make([]structs.PatronDisplay, 0),
    }

    // for {key}, {value} := range {list}
    for i, patron := range patrons.Data {
        attr := patron.Attributes

        // is an active patron that pays $25 or more
        if attr.PatronStatus == "active_patron" /*&& attr.WillPayAmountCents >= 2500*/ {
            userId := patron.Relationships.User.Data.Id
            imageUrl := patrons.Included[i].Attributes.ImageUrl

            newCache.Patrons = append(newCache.Patrons, structs.PatronDisplay{
                Id:       userId,
                Name:     attr.FullName,
                ImageUrl: imageUrl,
            })
        }
    }

    log.Println("Found", len(newCache.Patrons), "patrons")

    cache.PatronCache = newCache
}

func StartUpdatePatronTimer() {
    go UpdatePatrons()

    // tick every 24 hours to update the patrons
    ticker := time.NewTicker(24 * time.Hour)

    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)

    go func() {
        for {
            select {
            case <- ticker.C:
                UpdatePatrons()
            case <- sigint:
                ticker.Stop()
                return
            }
        }
    }()
}

func LoadPatronCredentials() {
    if _, err := os.Stat("/storage/oengus-patreon/patreon-credentials.json"); os.IsNotExist(err) {
        // fetch credentials
        log.Fatal("Failed to load credentials file at \"/storage/oengus-patreon/patreon-credentials.json\"")
    }

    // Open our jsonFile
    jsonFile, err := os.Open("/storage/oengus-patreon/patreon-credentials.json")
    // if we os.Open returns an error then handle it
    if err != nil {
        log.Fatal(err)
    }

    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal(byteValue, &cache.Tokens)
}
