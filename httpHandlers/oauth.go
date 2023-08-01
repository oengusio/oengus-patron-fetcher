package httpHandlers

import (
	"encoding/json"
	"net/http"
	"oenugs-patreon/patreon"
	"os"
	"strings"
)

// https://play.golang.org/p/Qg_uv_inCek
// contains checks if a string is present in a slice
func arrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

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

	origin := r.Header.Get("Origin")

	if origin == "" {
		origin = "https://oengus.io"
	}

	allowedOrigins := strings.Split(os.Getenv("OENGUS_BASE"), ",")

	if !arrayContains(allowedOrigins, origin) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Origin is not allowed"}`))
		return
	}

	token, err := patreon.Oauth2FetchToken(code, origin)

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
	data, _ := json.Marshal(user)

	w.Write(data)
}
