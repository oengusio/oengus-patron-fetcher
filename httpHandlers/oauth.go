package httpHandlers

import (
    "fmt"
    "log"
    "net/http"
    "oenugs-patreon/patreon"
)

func OauthAuthorize(w http.ResponseWriter, r *http.Request) {
    token, err := patreon.Oauth2FetchToken(r.URL.Query().Get("code"))

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(fmt.Sprintf(
            `{"error": "%s"}`,
            err.Error(),
        )))
        return
    }

    log.Println(token)
}
