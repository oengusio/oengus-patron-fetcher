package httpHandlers

import (
    "encoding/json"
    "net/http"
    "oenugs-patreon/cache"
)

func Patrons(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte(`{"error": "Method not allowed"}`))
        return
    }

    theJson, _ := json.Marshal(cache.PatronCache)

    w.Write(theJson)
}
