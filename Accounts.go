package main

import (
	"github.com/labstack/echo"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"fmt"
	"passwordresetrequest"
	//"strconv"
)

type Account struct {
	ID			int64	`json:id`
	UserID		int64	`json:userid`
	Email		string	`json:email`
	Password	string	`json:password`
	Locks		int64	`json:locks`
	Banned		bool	`json:banned`
}

type CreationRequest struct {
	ID			int64	`json:id`
	UserID		int64	`json:userid`
	Email		string	`json:email`
	DisplayName	string	`json:display_name`
	Password 	string	`json:password`
	Age			int64	`json:age`
	AccountID	int64	`json:accountid`
	StartTime	int64	`json:start_time`
	EndTime		int64	`json:end_time`
	Status		string	`json:status`
}

type RequestResult struct {
	Success		bool	`json: success`
	Message		string	`json: message`
	Result 		interface{} `json: result`
}

func getResetStatus(c echo.Context) error {
	request := new(PasswordResetRequest)
	if err = c.Bind(request); err != nil {
		return nil
	}

	session := pool.NewSession(nil)
	

	if err = session.Select("*").From("password_request").Where("id = ?", request.ID).LoadOne(&request); err != nil {
		return c.JSON(http.StatusNotFound, &RequestResult{ Success: false, Message: fmt.Sprintf("Password reset request with the ID %d does not exist.", request.ID) })
	}

	return c.JSON(http.StatusOK, &RequestResult{ Success: true, Result: request })
}

func resetPassword(c echo.Context) error {
	request := new(PasswordResetRequest)
	if err = c.Bind(request); err != nil {
		return nil
	}

	request.StartTime = timestamp()

	i, _ := strconv.ParseInt(c.PostForm("accountid"), 10, 64)
	accountid := int64(i)
	newPassword := c.PostForm("new_password")
	request := &PasswordResetRequest{ AccountID: accountid, NewPassword: newPassword, StartTime: timestamp() }

	session := pool.NewSession(nil)
	session.InsertInto("password_request").Columns("AccountID", "NewPassword").Record(&request).Exec()

	return c.JSON(http.StatusOK, request)
}

/*func getCreateStatus(c echo.Context) error {
	var requests []CreationRequest
	requestid := c.Param("id")
	
	session := pool.NewSession(nil)
	session.Select("*").From("creation_request").Where("id = ?", requestid).Load(&requests)

	return c.JSON(http.StatusOK, requests)
}

func createAccount(c echo.Context) error {
	email := c.PostForm("email")
	password := c.PostForm("password")
	displayName := c.PostForm("display_name")
	i, _ := strconv.ParseInt(c.PostForm("age"), 10, 64)
	age := int64(i)
	request := &CreationRequest{ Email: email, Password: password, DisplayName: displayName, Age: age, StartTime: timestamp() }

	session := pool.NewSession(nil)
	session.InsertInto("creation_request").Columns("email", "password").Record(&request).Exec()

	return c.JSON(http.StatusOK, request)
}

func addAccount(c echo.Context) error {
	email := c.Param("email")
	password := c.Param("password")
	account := &Account{ Email: email, Password: password }

	session := pool.NewSession(nil)
	session.InsertInto("accounts").Columns("email", "password").Record(&account).Exec()

	return c.JSON(http.StatusOK, account)
}

func getAccount(c echo.Context) error {
	var accounts []Account
	id := c.Param("id")

	session := pool.NewSession(nil)
	session.Select("*").From("accounts").Where("id = ?", id).Load(&accounts)

	return c.JSON(http.StatusOK, accounts)
}

func getAccounts(c echo.Context) error {
	var accounts []Account

	session := pool.NewSession(nil)
	session.Select("*").From("accounts").Load(&accounts)

	return c.JSON(http.StatusOK, accounts)
}

func deleteAccount(c echo.Context) error {
	id := c.Param("id")

	session := pool.NewSession(nil)
	session.DeleteFrom("accounts").Where("id = ?", id)

	return c.JSON(http.StatusOK, gin.H{ "success": true })
}

func updateAccount(c echo.Context) error {
	id := c.Param("id")

	session := pool.NewSession(nil)
	session.Update("accounts").Set("email = ?", c.Param("email")).Set("password = ?", c.Param("password")).Where("id = ?", id)

	return c.JSON(http.StatusOK, gin.H{ "success": true })
}*/