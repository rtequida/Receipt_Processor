package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rtequida/Receipt_Processor/api"
	"github.com/rtequida/Receipt_Processor/handlers"
)

func main() {
	e := echo.New()
	handler := handlers.NewReceiptHandler()
	api.RegisterHandlers(e, handler)
	e.Logger.Fatal(e.Start(":8080"))
}
