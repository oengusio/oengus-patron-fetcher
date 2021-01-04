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
    log.Println("Updating patrons")

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
        if attr.PatronStatus == "active_patron" /*&& attr.WillPayAmountCents >= 2500*/ {
            newCache.Patrons = append(newCache.Patrons, PatronDisplay{
                Id: patron.Relationships.User.Data.Id,
                Name: attr.FullName,
            })
        }
    }

    log.Println("Found", len(newCache.Patrons), "patrons")

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

func logRequestHandler(h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers",
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

        // Stop here if its Preflighted OPTIONS request
        if r.Method == "OPTIONS" {
            return
        }

        // call the original http.Handler we're wrapping
        h.ServeHTTP(w, r)
    }

    // http.HandlerFunc wraps a function so that it
    // implements http.Handler interface
    return http.HandlerFunc(fn)
}

func main() {
    InitApp()

    mux := http.NewServeMux()

    // TODO: Turn these into nice functions :)
    mux.HandleFunc("/patrons", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        theJson, _ := json.Marshal(patreonCache)

        w.Write(theJson)
    })

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // prevent all routes from catching on "/"
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }

        w.Header().Set("Content-Type", "application/json")

        // gotta love indentation stuff with this :D
        w.Write([]byte(`{
    "routes": {
        "/patrons": "Fetches the oengus patrons that are on the $25 tier"
    }
}`))
    })

    var handler http.Handler = mux
    // wrap mux with our logger. this will
    handler = logRequestHandler(handler)

    server := &http.Server{
        Addr: ":9000",
        ReadTimeout:  120 * time.Second,
        WriteTimeout: 120 * time.Second,
        IdleTimeout:  120 * time.Second, // introduced in Go 1.8
        Handler:      handler,
    }

    log.Println("Listening to port 9000")
    log.Fatal(server.ListenAndServe())
}
