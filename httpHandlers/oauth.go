package httpHandlers

import (
    "fmt"
    "log"
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
        w.Write([]byte(fmt.Sprintf(
            `{"error": "%s"}`,
            err.Error(),
        )))
        return
    }

    log.Println(token)

    user, fetchErr := patreon.Oauth2FetchUser(token)

    if fetchErr != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(fmt.Sprintf(
            `{"error": "%s"}`,
            fetchErr.Error(),
        )))
        return
    }

    log.Println(user)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"error": "boop"}`))
}
