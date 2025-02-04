package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rtequida/Receipt_Processor/api"
	"github.com/stretchr/testify/assert"
)

func TestPostReceiptsProcess(t *testing.T) {
	e := echo.New()
	handler := NewReceiptHandler()

	JSON_Receipt_Valid := `{
		"retailer": "Test Market",
		"purchaseDate": "2086-10-10",
		"purchaseTime": "14:14",
		"items": [
			{
			"shortDescription": "Gatorade",
			"price": "2.25"
			}
		],
		"total": "2.25"
	}`
	JSON_Receipt_Invalid := `{
		"retailer": "Test Market",
		"purchaseDate": "2086-Feb-10",
		"purchaseTime": "14:14",
		"items": [
			{
			"shortDescription": "Gatorade",
			"price": "2.25"
			}
		],
		"total": "2.25"
	}`

	t.Run("Valid Receipt", func(t *testing.T) {
		// Creating a new mock request to be sent to the echo server
		req := httptest.NewRequest(http.MethodPost, "/receipts/process", strings.NewReader(JSON_Receipt_Valid))
		// Mocking the headers of the request
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		// A mock response recorder to handle the response from the handler
		rec := httptest.NewRecorder()
		// Creating a new context for echo
		ctx := e.NewContext(req, rec)
		// Calls the POST handler and records response to rec
		err := handler.PostReceiptsProcess(ctx)
		// Makes sure the response had no errors
		assert.NoError(t, err)
		// Make sure the response has the correct status code
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]string
		// Creates a go map from the JSON response
		json.Unmarshal(rec.Body.Bytes(), &response)
		// Make sure the id field is included in the response
		assert.Contains(t, response, "id")
	})

	t.Run("Invalid Receipt", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/receipts/process", strings.NewReader(JSON_Receipt_Invalid))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		err := handler.PostReceiptsProcess(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Contains(t, "The receipt is invalid.", response["error"])
	})
}

func TestGetReceiptsIdPoints(t *testing.T) {
	e := echo.New()
	handler := NewReceiptHandler()

	receipt := api.Receipt{
		Retailer:     "Test Market",
		PurchaseDate: openapi_types.Date{Time: time.Date(2001, 05, 15, 0, 0, 0, 0, time.UTC)},
		PurchaseTime: "14:14",
		Items:        []api.Item{{ShortDescription: "lemons", Price: "5.00"}},
		Total:        "5.00",
	}
	id := "valid-id"
	handler.receipts[id] = receipt

	t.Run("Valid Id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/receipts/"+id+"/points", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		err := handler.GetReceiptsIdPoints(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Contains(t, response, "points")
	})

	t.Run("Invalid Receipt", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/receipts/not-a-real-id/points", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		err := handler.GetReceiptsIdPoints(ctx, "not-a-real-id")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Contains(t, "No receipt found for that ID.", response["error"])
	})
}
