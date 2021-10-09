package main

import (
    "log"
    "net/http"
    "oenugs-patreon/cache"
    "oenugs-patreon/httpHandlers"
    "oenugs-patreon/structs"
    "os"
    "time"
)

func InitApp() {
    // Store default patron array
    cache.PatronCache = structs.PatronOutput{
        Patrons: make([]structs.PatronDisplay, 0),
    }

    // only load the patrons if we want to
    if os.Getenv("DISABLE_CLOCK") == "false" {
        log.Println("Starting clock")

        // Load the stored tokens
        LoadPatronCredentials()

        // run the update func in a goroutine
        StartUpdatePatronTimer()
    }
}

func logRequestHandler(h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "*")

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

    mux.HandleFunc("/sync", httpHandlers.OauthAuthorize)
    mux.HandleFunc("/patrons", httpHandlers.Patrons)
    mux.HandleFunc("/webhook", httpHandlers.Webhook)
    mux.HandleFunc("/", httpHandlers.Root)

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
