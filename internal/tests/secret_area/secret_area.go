package secret_area

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"seal/internal/domain/secret_area"
	"seal/internal/tests/data"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData *data.TestData
var created secret_area.SecretArea

func Run(t *testing.T, data *data.TestData) {
	testData = data

	add(t)
	update(t)
	list(t)
	get(t)
	del(t)
}

func add(t *testing.T) {
	payload := fmt.Sprintf(`{"title":"test_%d","area":[[56.173023,35.777077],[57.16499,39.68821]],"description": "createdForTest"}`,
		testData.TimeStamp)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/secret-area", strings.NewReader(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&created), nil)
	assert.NotEmpty(t, created.Title)
	assert.NotEmpty(t, created.Description)
	assert.NotEmpty(t, created.Area)
}

func update(t *testing.T) {
	payload := fmt.Sprintf(`{"title":"test_update_%d","area":[[56.173323,35.777027],[57.16499,39.68821]],"description": "createdForTest_update"}`,
		testData.TimeStamp)

	url := fmt.Sprintf("/api/v1/secret-area/%d", created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&created), nil)
}

func list(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/secret-area?find=%s`, created.Title)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	var listResp struct {
		RecordsFiltered int                      `json:"records_filtered"`
		RecordsTotal    int                      `json:"records_total"`
		Data            []secret_area.SecretArea `json:"data"`
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&listResp), nil)
	assert.NotEqual(t, listResp.RecordsFiltered, 0)
	assert.NotEmpty(t, listResp.Data)
}

func get(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/secret-area/%d`, created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func del(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/secret-area/%d`, created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
