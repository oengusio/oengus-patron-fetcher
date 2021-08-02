package httpHandlers

import (
    "encoding/json"
    "net/http"
    "oenugs-patreon/patreon"
)

func OauthAuthorize(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte(`{"error": "Method not allowed"}`))
        return
    }

    code := r.URL.Query().Get("code")

    if code == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "missing code in request"}`))
        return
    }

    token, err := patreon.Oauth2FetchToken(code)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    user, fetchErr := patreon.Oauth2FetchUser(token)

    if fetchErr != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(fetchErr.Error()))
        return
    }

    w.WriteHeader(http.StatusOK)
    data,_ := json.Marshal(user)

    w.Write(data)
}
