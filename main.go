package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
)

var tokens PatreonTokens
var patreonCache PatronOutput

func UpdatePatrons() {
    patrons, err := FetchPatrons(tokens)

    // 401 response, refresh the tokens
    if err != nil && err.Error() == "StatusUnauthorized" {
        newTokens, _ := RefreshToken(tokens)

        tokens = newTokens

        patrons, err = FetchPatrons(tokens)
    }

    if err != nil {
        log.Println(err)
        return
    }

    newCache := PatronOutput{
        Patrons: make([]PatronDisplay, 0),
    }

    // for {key}, {value} := range {list}
    for _, patron := range patrons.Data {
        attr := patron.Attributes

        // is an active patron that pays $25 or more
        if attr.PatronStatus == "active_patron" && attr.WillPayAmountCents >= 2500 {
            newCache.Patrons = append(newCache.Patrons, PatronDisplay{
                Id: patron.Relationships.User.Data.Id,
                Name: attr.FullName,
            })
        }
    }

    patreonCache = newCache
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
    // create cache folder
    if _, err := os.Stat("cache"); os.IsNotExist(err) {
        os.Mkdir("cache", 0700)
    }

    if _, err := os.Stat("cache/patreon-credentials.json"); os.IsNotExist(err) {
        // fetch credentials
        log.Fatal("Failed to load credentials file, please place it in the cache folder")
    }

    // Open our jsonFile
    jsonFile, err := os.Open("cache/patreon-credentials.json")
    // if we os.Open returns an error then handle it
    if err != nil {
        log.Fatal(err)
    }

    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal(byteValue, &tokens)
}

func InitApp() {
    // Load the stored tokens
    LoadPatronCredentials()

    // Store default patron array
    patreonCache = PatronOutput{
        Patrons: make([]PatronDisplay, 0),
    }

    // run the update func in a goroutine
    StartUpdatePatronTimer()
}

func main() {
    InitApp()

    mux := http.NewServeMux()

    // Using /patrons here becuase "/" would catch all routes sent to it
    mux.HandleFunc("/patrons", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        theJson, _ := json.Marshal(patreonCache)

        w.Write(theJson)
    })

    log.Println("Listening to port 9000")

    log.Fatal(http.ListenAndServe(":9000", mux))
}
