package main

import (
	"github.com/chtavares592/client-auth/handlers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/", handlers.HandlerSlash)
	e.GET("/callback", handlers.HandlerCallback)

	e.Start(":1323")

}
