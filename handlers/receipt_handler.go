// Package handlers implements the handlers that connect the api routes to the backend functionality.
package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rtequida/Receipt_Processor/api"
	"github.com/rtequida/Receipt_Processor/app"
)

// A ReceiptHanlder an ID to a Receipt for in-memory storage of Receipts.
type ReceiptHandler struct {
	receipts map[string]api.Receipt
}

// NewReceiptHandler creates a new instance of a ReceiptHandler and returns the pointer to it.
func NewReceiptHandler() *ReceiptHandler {
	return &ReceiptHandler{receipts: make(map[string]api.Receipt)}
}

// PostReceiptsProcess handles the POST method from the API and ties it to the backend to generate
// and store an ID from the given receipt. This function also validates the given payload to ensure
// all required fields are included and meet the schema patterns.
func (h *ReceiptHandler) PostReceiptsProcess(ctx echo.Context) error {
	var receipt api.Receipt
	if err := ctx.Bind(&receipt); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "The receipt is invalid."})
	}
	fmt.Println("Receipt: ", receipt)
	if flag, err_message := app.ValidateReceipt(receipt); !flag {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err_message})
	}
	id := app.GenerateId(receipt)
	h.receipts[id] = receipt
	return ctx.JSON(http.StatusOK, map[string]string{"id": id})
}

// GetReceiptsIdPoints handles the GET method from the API and ties it to the backend to get the
// overall score of the given receipt. This function also validates the given id to ensure
// the id meets the schema pattern.
func (h *ReceiptHandler) GetReceiptsIdPoints(ctx echo.Context, id string) error {
	if flag, err_message := app.ValidateID(id); !flag {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err_message})
	}
	receipt, ok := h.receipts[id]
	if ok {
		points := app.GetPoints(receipt)
		return ctx.JSON(http.StatusOK, map[string]int{"points": points})
	}
	return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "No receipt found for that ID."})
}
