package httpHandlers

import (
    "net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
    // prevent all routes from catching on "/"
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte(`{"error": "Method not allowed"}`))
        return
    }

    w.Header().Set("Content-Type", "application/json")

    // gotta love indentation stuff with this :D
    w.Write([]byte(`{
    "routes": {
        "/patrons": "Fetches the oengus patrons that are on the $25 tier",
        "/webhook": "Patreon webhook callback",
        "/sync": "Sync account via oauth2"
    }
}`))
}
