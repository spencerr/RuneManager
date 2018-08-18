package main

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"fmt"
)

var (
	upgrader = websocket.Upgrader{}
)

func main() {
	setupDatabase()

	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH,    echo.POST, echo.DELETE},
	}))
	

	router.GET("/", func (c echo.Context) error {
		return c.String(http.StatusOK, "")
	})

	router.GET("/ws", socketHandler)

	accounts := router.Group("/accounts")
	/*accounts.GET("/", getAccounts).Name = "get-all-accounts"
	accounts.POST("/", addAccount).Name = "add-account"

	accounts.GET("/:id", getAccount).Name = "get-account"
	accounts.PATCH("/:id", updateAccount).Name = "update-account"
	accounts.DELETE("/:id", deleteAccount).Name = "delete-account"

	accounts.POST("/create", createAccount)
	accounts.GET("/create/:id", getCreateStatus)

	accounts.POST("/resetpw", resetPassword)*/
	accounts.GET("/resetpw/:id", getResetStatus)

	router.Logger.Fatal(router.Start(":8080"))
}

func socketHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	for {
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

