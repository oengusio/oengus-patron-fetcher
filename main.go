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
    log.Println("Fetching patrons!")
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
    } else {
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

    log.Println(tokens.AccessToken)
}

func InitApp() {
    LoadPatronCredentials()

    // Store default patron array
    patreonCache = PatronOutput{
        Patrons: make([]PatronDisplay, 0),
    }

    // Load the stored tokens

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
