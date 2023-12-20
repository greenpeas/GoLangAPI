package shipping

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"seal/internal/domain/shipping"
	"seal/internal/tests/data"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData *data.TestData
var created shipping.Shipping

func Run(t *testing.T, data *data.TestData) {
	testData = data

	add(t)
	start(t)
	wrongStart(t)
	end(t)
	wrongEnd(t)
	update(t)
	list(t)
	get(t)
	getRoute(t)
	getTelemetry(t)
	del(t)
}

func add(t *testing.T) {
	numbers := fmt.Sprintf("%d", testData.TimeStamp)
	val := numbers[len(numbers)-6:]

	payload := fmt.Sprintf(`{"custom_number": "00%s", "create_date": "%s", "number": %d, "transport": %d, 
"route": %d, "seal": %d}`, val, val, testData.TimeStamp,
		testData.Transport.Id, testData.Route.Id, testData.Seal.Id)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shipping", strings.NewReader(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&created), nil)
	assert.NotEmpty(t, created.Seal.SerialNumber)
}

func start(t *testing.T) {
	url := fmt.Sprintf("/api/v1/shipping/%d/start", created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func wrongStart(t *testing.T) {
	url := fmt.Sprintf("/api/v1/shipping/%d/start", created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func end(t *testing.T) {
	url := fmt.Sprintf("/api/v1/shipping/%d/end", created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func wrongEnd(t *testing.T) {
	url := fmt.Sprintf("/api/v1/shipping/%d/end", created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func update(t *testing.T) {
	numbers := fmt.Sprintf("%d", testData.TimeStamp)
	val := numbers[len(numbers)-6:]

	payload := fmt.Sprintf(`{"custom_number": "00%s"}`, val)

	url := fmt.Sprintf("/api/v1/shipping/%d", created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(payload))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&created), nil)
}

func list(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/shipping?find=%s`, created.CustomNumber)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	var listResp struct {
		RecordsFiltered int                 `json:"records_filtered"`
		RecordsTotal    int                 `json:"records_total"`
		Data            []shipping.Shipping `json:"data"`
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&listResp), nil)
	assert.NotEqual(t, listResp.RecordsFiltered, 0)
	assert.NotEmpty(t, listResp.Data)
}

func get(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/shipping/%d`, created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func getRoute(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/shipping/%d/route`, created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func getTelemetry(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/shipping/%d/telemetry?from=%s&limit=500`, created.Id, "2023-02-16T15%3A59%3A37.398774%2B03%3A00")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func del(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/shipping/%d`, created.Id)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
