package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"github.com/fatih/structs"
)

var (
	router *echo.Echo
)

func main() {
	loadConfig()
	setupDatabase()

	router = echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{ "*" },
		AllowMethods: []string{ echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE },
	}))
	

	router.File("/", "public/index.html")

	router.GET("/ws", socketHandler)

	accounts := router.Group("/accounts")
	accounts.GET("/all", accountHandler).Name = "get-all-accounts"
	accounts.POST("/", accountHandler).Name = "add-account"

	accounts.GET("/:id", accountHandler).Name = "get-account"
	accounts.PATCH("/:id", accountHandler).Name = "update-account"
	accounts.DELETE("/:id", accountHandler).Name = "delete-account"
	accounts.GET("/:id/lock", accountHandler).Name = "lock-account"
	accounts.GET("/:id/ban", accountHandler).Name = "ban-account"

	/*accounts.POST("/create", createAccount).Name = "create-account"
	accounts.GET("/create/:id", getCreateStatus).Name = "get-create-account-status"

	accounts.POST("/resetpw", resetPassword).Name = "reset-password"
	accounts.GET("/resetpw/:id", getResetStatus).Name = "get-reset-password-status"

	servers := router.Group("/servers")
	servers.GET("/all", getServers).Name = "get-all-servers"
	servers.POST("/", addServer).Name = "add-server"

	servers.GET("/:id", getServer).Name = "get-server"
	servers.PATCH("/:id", updateServer).Name = "update-server"
	servers.DELETE("/:id", deleteAccount).Name = "delete-server"

	servers.GET("/:id/clients", getClientsForServer).Name = "get-all-clients-for-server"

	clients := router.Group("/clients")
	clients.POST("/all", addClient).Name = "add-client"

	clients.GET("/:id", getClient).Name = "get-client"
	clients.PATCH("/:id", updateClient).Name = "update-client"
	clients.DELETE("/:id", deleteClient).Name = "delete-client"*/

	go hub.run()
	router.Logger.Fatal(router.Start(":8080"))
}

func accountHandler(c echo.Context) error {
	arguments := structs.Map(BindAccount(c))
	arguments["ApiKey"] = "asdf1234"
	return c.JSON(http.StatusOK, functions[getRouteName(c)](&APIRequest{ ApiArguments: arguments }))
}

func getRouteName(c echo.Context) string {
	for _, route := range router.Routes() {
		if route.Path == c.Path() && route.Method == c.Request().Method {
			return route.Name
		}
	}

	return ""
}