package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"seal/internal/tests/data"
	"seal/internal/transport"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData *data.TestData
var jwt transport.JwtWithRefresh

func Run(t *testing.T, data *data.TestData) transport.JwtWithRefresh {
	testData = data

	wrongSignIn(t)
	signIn(t)
	refreshToken(t)

	return jwt
}

func wrongSignIn(t *testing.T) {
	payload := `{"login": "not_existing_user", "password": "not_existing_password"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", strings.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func signIn(t *testing.T) {
	payload := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testData.App.Cfg.Test.User.Login, testData.App.Cfg.Test.User.Password)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", strings.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&jwt), nil)
}

func refreshToken(t *testing.T) {
	payload := fmt.Sprintf(`{"refresh_token": "%s"}`, jwt.RefreshToken)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&jwt), nil)
}
