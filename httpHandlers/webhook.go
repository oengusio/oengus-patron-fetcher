package httpHandlers

import (
    "crypto/hmac"
    "crypto/md5"
    "io/ioutil"
    "net/http"
    "os"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte(`{"error": "Method not allowed"}`))
        return
    }

    secret := r.Header.Get("X-Patreon-Signature")

    if secret == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "Missing X-Patreon-Signature header"}`))
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "Could not read body"}`))
        return
    }

    hash := hmac.New(md5.New, []byte(os.Getenv("PATREON_WEBHOOK_SECRET")))
    hash.Write(body)

    if !hmac.Equal([]byte(secret), hash.Sum(nil)) {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "X-Patreon-Signature header is not valid"}`))
        return
    }

    event := r.Header.Get("X-Patreon-Event")

    switch event {
    case "members:pledge:create":
        // new patron
        break
    case "members:pledge:update":
        // Updated patron
        break
    case "members:pledge:delete":
        // Remove perks
        break
    default:
        w.Write([]byte(`{"error": "event not programmed"}`))
        break
    }
}
