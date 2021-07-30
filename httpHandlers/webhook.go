package httpHandlers

import (
    "crypto/hmac"
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "oenugs-patreon/sql"
    "oenugs-patreon/structs"
    "os"
    "strings"
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

    calculated := hex.EncodeToString(hash.Sum(nil))

    if calculated != secret {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "X-Patreon-Signature header is not valid"}`))
        return
    }

    //log.Println(string(body))

    event := r.Header.Get("X-Patreon-Event")

    switch event {
    case "members:pledge:create":
        w.WriteHeader(http.StatusOK)
        // new patron
        var pledge structs.WebhookPledge
        json.Unmarshal(body, &pledge)

        go addPledge(pledge)
        break
    case "members:pledge:update":
        w.WriteHeader(http.StatusOK)
        // Updated patron
        var pledge structs.WebhookPledge
        json.Unmarshal(body, &pledge)

        go updatePledge(pledge)
        break
    case "members:pledge:delete":
        w.WriteHeader(http.StatusOK)
        // Remove perks
        var pledge structs.WebhookPledge
        json.Unmarshal(body, &pledge)

        go deletePledge(pledge)
        break
    default:
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "event not programmed"}`))
        break
    }
}

func addPledge(pledge structs.WebhookPledge) {
    status := parseStatus(pledge.Data.Attributes.PatronStatus)

    if status == "" {
        deletePledge(pledge)
        return
    }

    userId := pledge.Data.Relationships.User.Data.Id
    payAmount := pledge.Data.Attributes.PledgeAmountCents

    go sql.InsertMember(userId, status, payAmount)
}

func updatePledge(pledge structs.WebhookPledge) {
    status := parseStatus(pledge.Data.Attributes.PatronStatus)

    if status == "" {
        deletePledge(pledge)
        return
    }

    userId := pledge.Data.Relationships.User.Data.Id
    payAmount := pledge.Data.Attributes.PledgeAmountCents

    go sql.UpdateMember(userId, status, payAmount)
}

func deletePledge(pledge structs.WebhookPledge) {
    userId := pledge.Data.Relationships.User.Data.Id

    go sql.DeleteMember(userId)
}

func parseStatus(status string) string {
    return strings.ToUpper(status)
}
