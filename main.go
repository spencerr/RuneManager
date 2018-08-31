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

	router = echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{ "*" },
		AllowMethods: []string{ echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE },
	}))
	

	router.File("/", "public/index.html")
	router.GET("/ws", SocketHandler)

	accounts := router.Group("/accounts")
	accounts.GET("/all", AccountHandler).Name = "get-all-accounts"
	accounts.POST("/", AccountHandler).Name = "add-account"

	accounts.GET("/:id", AccountHandler).Name = "get-account"
	accounts.PATCH("/:id", AccountHandler).Name = "update-account"
	accounts.DELETE("/:id", AccountHandler).Name = "delete-account"
	accounts.GET("/:id/lock", AccountHandler).Name = "lock-account"
	accounts.GET("/:id/ban", AccountHandler).Name = "ban-account"

	accounts.POST("/create", AccountHandler).Name = "create-account"
	accounts.GET("/create/:id", AccountHandler).Name = "get-create-account-status"

	accounts.POST("/resetpw", AccountHandler).Name = "reset-password"
	accounts.GET("/resetpw/:id", AccountHandler).Name = "get-reset-password-status"

	servers := router.Group("/servers")
	servers.GET("/all", ServerHandler).Name = "get-all-servers"
	servers.POST("/", ServerHandler).Name = "add-server"

	servers.GET("/:id", ServerHandler).Name = "get-server"
	servers.PATCH("/:id", ServerHandler).Name = "update-server"
	servers.DELETE("/:id", ServerHandler).Name = "delete-server"

	servers.GET("/:id/clients", ServerHandler).Name = "get-all-clients-for-server"

	/*clients := router.Group("/clients")
	clients.POST("/all", addClient).Name = "add-client"

	clients.GET("/:id", getClient).Name = "get-client"
	clients.PATCH("/:id", updateClient).Name = "update-client"
	clients.DELETE("/:id", deleteClient).Name = "delete-client"*/

	go hub.run()
	router.Logger.Fatal(router.Start(":8080"))
}

func AccountHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Functions[GetRouteName(c)](&APIRequest{ ApiArguments: structs.Map(BindAccount(c)) }))
}

func ServerHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Functions[GetRouteName(c)](&APIRequest{ ApiArguments: structs.Map(BindServer(c)) }))
}

func GetRouteName(c echo.Context) string {
	for _, route := range router.Routes() {
		if route.Path == c.Path() && route.Method == c.Request().Method {
			return route.Name
		}
	}

	return "no-route"
}

func NoRoute(request *APIRequest) *APIResponse {
	return APIFail("no-route")
}