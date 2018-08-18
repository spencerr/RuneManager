package main

import (
	"github.com/labstack/echo"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"fmt"
	//"strconv"
)


func getResetStatus(c echo.Context) error {
	request := BindPasswordResetRequest(c)
	if (request == nil) {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: "Invalid password reset request id given." })
	}

	session := pool.NewSession(nil)
	if err := session.Select("*").From("password_reset_request").Where("id = ?", request.ID).LoadOne(&request); err != nil {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: fmt.Sprintf("Password reset request with the ID %d does not exist.", request.ID) })
	}

	return c.JSON(http.StatusOK, &RequestResult{ Success: true, Result: request })
}

func resetPassword(c echo.Context) error {
	request := BindPasswordResetRequest(c)
	if request == nil {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: "Unable to create a password reset request. Please ensure you are passing the required values." })
	}

	request.StartTime = timestamp()
	session := pool.NewSession(nil)
	session.InsertInto("password_reset_request").Columns("AccountID", "NewPassword").Record(&request).Exec()

	return c.JSON(http.StatusOK, request)
}

func getCreateStatus(c echo.Context) error {
	request := BindAccountCreationRequest(c)
	if request == nil {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: "Invalid account creation request id given." })
	}

	session := pool.NewSession(nil)
	if err := session.Select("*").From("account_creation_request").Where("id = ?", request.ID).LoadOne(&request); err != nil {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: fmt.Sprintf("Account creation request with the ID %d does not exist.", request.ID) })
	}

	return c.JSON(http.StatusOK, &RequestResult{ Success: true, Result: request })
}

func createAccount(c echo.Context) error {
	request := BindAccountCreationRequest(c)
	if request == nil {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: "Unable to create an account creation request. Please ensure you are passing the required values." })
	}

	request.StartTime = timestamp()
	session := pool.NewSession(nil)
	session.InsertInto("account_creation_request").Columns("email", "password", "start_time", "age").Record(&request).Exec()
	
	return c.JSON(http.StatusOK, &RequestResult{ Success: true, Result: request })
}

func addAccount(c echo.Context) error {
	account := BindAccount(c)

	session := pool.NewSession(nil)
	session.InsertInto("accounts").Columns("email", "password").Record(&account).Exec()

	return c.JSON(http.StatusOK, &RequestResult{ Success: true, Result: account })
}

func getAccount(c echo.Context) error {
	account := BindAccount(c)
	if account == nil {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: "Invalid account id given." })
	}

	session := pool.NewSession(nil)
	if err := session.Select("*").From("accounts").Where("id = ?", account.ID).LoadOne(&account); err != nil {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: fmt.Sprintf("Account with the ID %d does not exist.", account.ID) })
	}

	return c.JSON(http.StatusOK, &RequestResult{ Success: true, Result: account })
}

func getAccounts(c echo.Context) error {
	var accounts []Account

	session := pool.NewSession(nil)
	session.Select("*").From("accounts").Load(&accounts)

	return c.JSON(http.StatusOK, &RequestResult{ Success: true, Result: accounts })
}

func deleteAccount(c echo.Context) error {
	account := BindAccount(c)
	if account == nil {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: "Invalid account id given." })
	}

	success := false
	session := pool.NewSession(nil)
	if result, err := session.DeleteFrom("accounts").Where("id = ?", account.ID).Exec(); err != nil {
		if rows, err := result.RowsAffected(); err != nil {
			success = rows == 1
		}
	}

	return c.JSON(http.StatusOK, &RequestResult{ Success: success })
}

func updateAccount(c echo.Context) error {
	account := BindAccount(c)
	if account == nil {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: "Invalid account information given." })
	}

	session := pool.NewSession(nil)
	session.Update("accounts").Set("email = ?", account.Email).Set("password = ?", account.Password).Where("id = ?", account.ID).Exec()

	return c.JSON(http.StatusOK, &RequestResult{ Success: true, Result: account })
}