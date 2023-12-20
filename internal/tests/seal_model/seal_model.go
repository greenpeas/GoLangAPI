package seal_model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"seal/internal/domain/seal_model"
	"seal/internal/tests/data"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData *data.TestData

func Run(t *testing.T, data *data.TestData) {
	testData = data

	list(t)
}

func list(t *testing.T) {
	url := fmt.Sprintf(`/api/v1/seal-model?find=%s`, "Автомобильный")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testData.Jwt.Token))
	req.Header.Add("Content-Type", "application/json")

	testData.App.Router.ServeHTTP(w, req)

	var listResp struct {
		RecordsFiltered int                    `json:"records_filtered"`
		RecordsTotal    int                    `json:"records_total"`
		Data            []seal_model.SealModel `json:"data"`
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json.NewDecoder(w.Body).Decode(&listResp), nil)
	assert.NotEqual(t, listResp.RecordsFiltered, 0)
	assert.NotEmpty(t, listResp.Data)
}
