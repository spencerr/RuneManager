package main

import (
	//"github.com/labstack/echo"
	_ "github.com/jmoiron/sqlx"
	//"net/http"
	"fmt"
	"reflect"
)


/*func getResetStatus(c echo.Context) error {
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

func updateAccount(c echo.Context) error {
	account := BindAccount(c)
	if account == nil {
		return c.JSON(http.StatusOK, &RequestResult{ Success: false, Message: "Invalid account information given." })
	}

	session := pool.NewSession(nil)
	session.Update("accounts").Set("email = ?", account.Email).Set("password = ?", account.Password).Where("id = ?", account.ID).Exec()

	return c.JSON(http.StatusOK, &RequestResult{ Success: true, Result: account })
}*/

func getAccounts(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey" }); !ok {
		return &APIResponse{ Success: false, Result: err }
	}
	
	var accounts []Account
	qs := "SELECT Accounts.ID, Accounts.UserID, Accounts.Email, Accounts.Password, Accounts.Locks, Accounts.Banned FROM Accounts JOIN Users ON Users.ID = Accounts.UserID WHERE Users.ApiKey = :ApiKey"
	if err := SelectAll(&accounts, qs, request.ApiArguments); err != nil {
		debug(fmt.Sprintf("Unable to retrieve accounts. %s", err))
		return &APIResponse{ Success: false, Result: "select-fail" }
	}

	debug(fmt.Sprintf("Successfully retrieved accounts %+v", accounts))
	return &APIResponse{ Success: true, Result: accounts }
}

func addAccount(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey", "Email", "Password" }); !ok {
		return &APIResponse{ Success: false, Result: err }
	}

	qs := `INSERT INTO Accounts(UserID, Email, Password) VALUES(SELECT ID FROM Users WHERE ApiKey = :ApiKey, :Email, :Password)`
	result, err := Insert(qs, request.ApiArguments)
	if err != nil {
		debug(fmt.Sprintf("Unable to insert account. %s", err))
		return &APIResponse{ Success: false, Result: "insert-fail" }
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		debug(fmt.Sprintf("Unable to get insertid. %s", err))
		return &APIResponse{ Success: false, Result: "no-insert-id" }
	}

	return &APIResponse{ Success: true, Result: lastInsertId }
}

func getAccount(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return &APIResponse{ Success: false, Result: err }
	}


	var account Account
	qs := "SELECT Accounts.ID, Accounts.UserID, Accounts.Email, Accounts.Password, Accounts.Locks, Accounts.Banned FROM Accounts JOIN Users ON Users.ID = Accounts.UserID WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID"
	if err := SelectOne(&account, qs, request.ApiArguments); err != nil {
		debug(fmt.Sprintf("Unable to retrieve account. %s", err))
		return &APIResponse{ Success: false, Result: "select-fail" }
	}

	debug(fmt.Sprintf("Successfully retrieved account %+v", account))
	return &APIResponse{ Success: true, Result: account }
}

func updateAccount(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey", "ID" }, "UserID"); !ok {
		return &APIResponse{ Success: false, Result: err }
	}

	qs := `UPDATE Accounts JOIN Users ON Users.ID = Accounts.UserID SET`
	keys := ``

	elem := reflect.ValueOf(&Account{}).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)
		if tag, ok := field.Tag.Lookup("static"); ok && tag == "true" {
			continue
		}

		key := field.Name
		if _, ok := request.ApiArguments[key]; ok {
			if len(keys) != 0 {
				keys += `,`
			}
			keys += ` Accounts.` + key + ` = :` + key
		}
	}

	result, err := Insert(qs + keys + ` WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID`, request.ApiArguments)
	if err != nil {
		debug(fmt.Sprintf("Unable to update account. %s", err))
		return &APIResponse{ Success: false, Result: "update-fail" }
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		debug(fmt.Sprintf("Unable to get rows affected. %s", err))
		return &APIResponse{ Success: false, Result: "no-rows-affected" }
	}


	return &APIResponse{ Success: true, Result: rowsAffected }
}


func deleteAccount(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return &APIResponse{ Success: false, Result: err }
	}

	qs := `DELETE Accounts FROM Accounts JOIN Users ON Users.ID = Accounts.UserID WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID`
	result, err := Delete(qs, request.ApiArguments)
	if err != nil {
		debug(fmt.Sprintf("Unable to delete account. %s", err))
		return &APIResponse{ Success: false, Result: "delete-fail" }
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		debug(fmt.Sprintf("Unable to get rows affected. %s", err))
		return &APIResponse{ Success: false, Result: "no-rows-affected" }
	}

	return &APIResponse{ Success: true, Result: rowsAffected }
}

func lockAccount(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return &APIResponse{ Success: false, Result: err }
	}

	qs := `UPDATE Accounts JOIN Users ON Users.ID = Accounts.UserID SET Locks = Locks + 1 WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID`
	result, err := Update(qs, request.ApiArguments)
	if err != nil {
		debug(fmt.Sprintf("Unable to lock account. %s", err))
		return &APIResponse{ Success: false, Result: "update-fail" }
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		debug(fmt.Sprintf("Unable to get rows affected. %s", err))
		return &APIResponse{ Success: false, Result: "no-rows-affected" }
	}

	return &APIResponse{ Success: true, Result: rowsAffected }
}

func banAccount(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return &APIResponse{ Success: false, Result: err }
	}

	qs := `UPDATE Accounts JOIN Users ON Users.ID = Accounts.UserID SET Banned = TRUE WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID`
	result, err := Update(qs, request.ApiArguments)
	if err != nil {
		debug(fmt.Sprintf("Unable to ban account. %s", err))
		return &APIResponse{ Success: false, Result: "update-fail" }
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		debug(fmt.Sprintf("Unable to get rows affected. %s", err))
		return &APIResponse{ Success: false, Result: "no-rows-affected" }
	}

	return &APIResponse{ Success: true, Result: rowsAffected }
}

func getResetStatus(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return &APIResponse{ Success: false, Result: err }
	}

	var req PasswordResetRequest
	qs := `SELECT * FROM PasswordResetRequests JOIN Users ON Users.ID = PasswordResetRequests.UserID WHERE User.ApiKey = :ApiKey AND PasswordResetRequest.ID = :ID`
	if err := SelectOne(&req, qs, request.ApiArguments); err != nil {
		debug(fmt.Sprintf("Unable to retrieve password reset request. %s", err))
		return &APIResponse{ Success: false, Result: "select-fail" }
	}

	debug(fmt.Sprintf("Successfully retrieved password reset request %+v", req))
	return &APIResponse{ Success: true, Result: req }
}

func resetPassword(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey", "NewPassword", "ID" }); !ok {
		return &APIResponse{ Success: false, Result: err }
	}

	qs := `INSERT INTO PasswordResetRequests(UserID, AccountID, NewPassword, StartTime) VALUES(SELECT ID FROM Users WHERE ApiKey = :ApiKey, :ID, :NewPassword, :StartTime)`
	request.ApiArguments["StartTime"] = timestamp()

	result, err := Insert(qs, request.ApiArguments)
	if err != nil {
		debug(fmt.Sprintf("Unable to insert account. %s", err))
		return &APIResponse{ Success: false, Result: "insert-fail" }
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		debug(fmt.Sprintf("Unable to get insertid. %s", err))
		return &APIResponse{ Success: false, Result: "no-insert-id" }
	}

	return &APIResponse{ Success: true, Result: lastInsertId }
}

func getCreateStatus(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return &APIResponse{ Success: false, Result: err }
	}

	var req AccountCreationRequest
	qs := `SELECT * FROM AccountCreationRequests JOIN Users ON Users.ID = AccountCreationRequests.UserID WHERE User.ApiKey = :ApiKey AND AccountCreationRequests.ID = :ID`
	if err := SelectOne(&req, qs, request.ApiArguments); err != nil {
		debug(fmt.Sprintf("Unable to retrieve account creation request. %s", err))
		return &APIResponse{ Success: false, Result: "select-fail" }
	}

	debug(fmt.Sprintf("Successfully retrieved account creation request %+v", req))
	return &APIResponse{ Success: true, Result: req }
}

func createAccount(request *APIRequest) *APIResponse {
	if ok, err := validateRequest(request, []string { "ApiKey", "Email", "DisplayName", "Password", "Age" }); !ok {
		return &APIResponse{ Success: false, Result: err }
	}

	qs := `INSERT INTO AccountCreationRequests(UserID, Email, DisplayName, Password, Age, StartTime) VALUES(SELECT ID FROM Users WHERE ApiKey = :ApiKey, :Email, :DisplayName, :Password, :Age, :StartTime)`
	request.ApiArguments["StartTime"] = timestamp()

	result, err := Insert(qs, request.ApiArguments)
	if err != nil {
		debug(fmt.Sprintf("Unable to insert account creation request. %s", err))
		return &APIResponse{ Success: false, Result: "insert-fail" }
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		debug(fmt.Sprintf("Unable to get insertid. %s", err))
		return &APIResponse{ Success: false, Result: "no-insert-id" }
	}

	return &APIResponse{ Success: true, Result: lastInsertId }
}
