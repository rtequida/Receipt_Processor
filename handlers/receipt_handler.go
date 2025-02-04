package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rtequida/Receipt_Processor/api"
	"github.com/rtequida/Receipt_Processor/app"
)

type ReceiptHandler struct {
	receipts map[string]api.Receipt
}

func NewReceiptHandler() *ReceiptHandler {
	return &ReceiptHandler{receipts: make(map[string]api.Receipt)}
}

func (h *ReceiptHandler) PostReceiptsProcess(ctx echo.Context) error {
	var receipt api.Receipt
	if err := ctx.Bind(&receipt); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "The receipt is invalid."})
	}
	flag, err_message := app.ValidateReceipt(receipt)
	if !flag {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err_message})
	}
	id := app.GenerateId(receipt)
	h.receipts[id] = receipt
	return ctx.JSON(http.StatusOK, map[string]string{"id": id})
}

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
