package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"seal/internal/domain/transport"
	"seal/internal/tests/data"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData *data.TestData

var created transport.Transport

func Run(t *testing.T, data *data.TestData) transport.Transport {
	testData = data

	add(t)
	update(t)
	list(t)
	get(t)

	return created
}

func add(t *testing.T) {
	payload := fmt.Sprintf(`{"title": "transport_%d", "author": %d, "type": 1}`, testData.TimeStamp, testData.User.Id)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transport", strings.NewReader(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&created), nil)
}

func update(t *testing.T) {
	payload := fmt.Sprintf(`{"title": "t_%d"}`, testData.TimeStamp)

	url := fmt.Sprintf("/api/v1/transport/%d", created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&created), nil)
}

func list(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/transport?find=%s`, "")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	var listResp struct {
		RecordsFiltered int                   `json:"records_filtered"`
		RecordsTotal    int                   `json:"records_total"`
		Data            []transport.Transport `json:"data"`
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&listResp), nil)
	assert.NotEqual(t, listResp.RecordsFiltered, 0)
	assert.NotEmpty(t, listResp.Data)
}

func get(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/transport/%d`, created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Delete(t *testing.T, testData *data.TestData) {
	url := fmt.Sprintf(`/api/v1/transport/%d`, testData.Transport.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
