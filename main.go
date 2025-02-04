package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rtequida/Receipt_Processor/api"
	"github.com/rtequida/Receipt_Processor/handlers"
)

// Main initializes the echo instance that handles requests to the defined
// routes in generatedserver.go which are registered with their respective
// handlers and corresponding functions. Then, the echo server is started and listening for requests.
func main() {
	e := echo.New()
	handler := handlers.NewReceiptHandler()
	api.RegisterHandlers(e, handler)
	e.Logger.Fatal(e.Start(":8080"))
}
