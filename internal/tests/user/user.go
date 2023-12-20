package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"seal/internal/domain/user"
	"seal/internal/tests/data"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData *data.TestData
var created user.User

func Run(t *testing.T, data *data.TestData) user.User {
	testData = data

	add(t)
	update(t)
	list(t)
	get(t)

	return created
}

func add(t *testing.T) {
	payload := fmt.Sprintf(`{"login": "test_%d", "password": "createdForTest", "role": 2}`, testData.TimeStamp)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user", strings.NewReader(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&created), nil)
}

func update(t *testing.T) {
	payload := fmt.Sprintf(`{"login": "t_%d"}`, testData.TimeStamp)

	url := fmt.Sprintf("/api/v1/user/%d", created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&created), nil)
}

func list(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/user?find=%s`, created.Login)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	var listResp struct {
		RecordsFiltered int         `json:"records_filtered"`
		RecordsTotal    int         `json:"records_total"`
		Data            []user.User `json:"data"`
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&listResp), nil)
	assert.NotEqual(t, listResp.RecordsFiltered, 0)
	assert.NotEmpty(t, listResp.Data)
}

func get(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/user/%d`, created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Delete(t *testing.T, testData *data.TestData) {
	url := fmt.Sprintf(`/api/v1/user/%d`, testData.User.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
