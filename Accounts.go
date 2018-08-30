package main

import (
	_ "github.com/jmoiron/sqlx"
	"fmt"
	"reflect"
)

func GetAccounts(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey" }); !ok {
		return APIFail(err)
	}
	
	var accounts []Account
	qs := "SELECT Accounts.ID, Accounts.UserID, Accounts.Email, Accounts.Password, Accounts.Locks, Accounts.Banned FROM Accounts JOIN Users ON Users.ID = Accounts.UserID WHERE Users.ApiKey = :ApiKey"
	if err := SelectAll(&accounts, qs, request.ApiArguments); err != nil {
		PrintError(fmt.Sprintf("Unable to retrieve accounts. %s", err))
		return APIFail(SelectFail)
	}

	PrintSuccess(fmt.Sprintf("Successfully retrieved accounts %+v", accounts))
	return APISuccess(accounts)
}

func AddAccount(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey", "Email", "Password" }); !ok {
		return APIFail(err)
	}

	qs := `INSERT INTO Accounts(UserID, Email, Password) VALUES((SELECT ID FROM Users WHERE ApiKey = :ApiKey), :Email, :Password)`
	result, err := Insert(qs, request.ApiArguments)
	if err != nil {
		PrintError(fmt.Sprintf("Unable to insert account. %s", err))
		return APIFail(InsertFail)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		PrintError(fmt.Sprintf("Unable to get insertid. %s", err))
		return APIFail(NoInsertID)
	}

	return APISuccess(lastInsertId)
}

func GetAccount(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return APIFail(err)
	}

	var account Account
	qs := "SELECT Accounts.ID, Accounts.UserID, Accounts.Email, Accounts.Password, Accounts.Locks, Accounts.Banned FROM Accounts JOIN Users ON Users.ID = Accounts.UserID WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID"
	if err := SelectOne(&account, qs, request.ApiArguments); err != nil {
		PrintError(fmt.Sprintf("Unable to retrieve account. %s", err))
		return APIFail(SelectFail)
	}

	PrintSuccess(fmt.Sprintf("Successfully retrieved account %+v", account))
	return APISuccess(account)
}

func UpdateAccount(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey", "ID" }, "UserID"); !ok {
		return APIFail(err)
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
		PrintError(fmt.Sprintf("Unable to update account. %s", err))
		return APIFail(UpdateFail)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		PrintError(fmt.Sprintf("Unable to get rows affected. %s", err))
		return APIFail(NoRowsAffected)
	}


	return APISuccess(rowsAffected)
}


func DeleteAccount(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return APIFail(err)
	}

	qs := `DELETE Accounts FROM Accounts JOIN Users ON Users.ID = Accounts.UserID WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID`
	result, err := Delete(qs, request.ApiArguments)
	if err != nil {
		PrintError(fmt.Sprintf("Unable to delete account. %s", err))
		return APIFail(DeleteFail)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		PrintError(fmt.Sprintf("Unable to get rows affected. %s", err))
		return APIFail(NoRowsAffected)
	}

	return APISuccess(rowsAffected)
}

func LockAccount(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return APIFail(err)
	}

	qs := `UPDATE Accounts JOIN Users ON Users.ID = Accounts.UserID SET Locks = Locks + 1 WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID`
	result, err := Update(qs, request.ApiArguments)
	if err != nil {
		PrintError(fmt.Sprintf("Unable to lock account. %s", err))
		return APIFail(UpdateFail)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		PrintError(fmt.Sprintf("Unable to get rows affected. %s", err))
		return APIFail(NoRowsAffected)
	}

	return APISuccess(rowsAffected)
}

func BanAccount(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return APIFail(err)
	}

	qs := `UPDATE Accounts JOIN Users ON Users.ID = Accounts.UserID SET Banned = TRUE WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID`
	result, err := Update(qs, request.ApiArguments)
	if err != nil {
		PrintError(fmt.Sprintf("Unable to ban account. %s", err))
		return APIFail(UpdateFail)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		PrintError(fmt.Sprintf("Unable to get rows affected. %s", err))
		return APIFail(NoRowsAffected)
	}

	return APISuccess(rowsAffected)
}

func GetResetPasswordStatus(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return APIFail(err)
	}

	var req PasswordResetRequest
	qs := `SELECT * FROM PasswordResetRequests JOIN Users ON Users.ID = PasswordResetRequests.UserID WHERE User.ApiKey = :ApiKey AND PasswordResetRequest.ID = :ID`
	if err := SelectOne(&req, qs, request.ApiArguments); err != nil {
		PrintError(fmt.Sprintf("Unable to retrieve password reset request. %s", err))
		return APIFail(SelectFail)
	}

	PrintSuccess(fmt.Sprintf("Successfully retrieved password reset request %+v", req))
	return APISuccess(req)
}

func ResetAccountPassword(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey", "NewPassword", "ID" }); !ok {
		return APIFail(err)
	}

	qs := `INSERT INTO PasswordResetRequests(UserID, AccountID, NewPassword, StartTime) VALUES((SELECT ID FROM Users WHERE ApiKey = :ApiKey), :ID, :NewPassword, :StartTime)`
	request.ApiArguments["StartTime"] = Timestamp()

	result, err := Insert(qs, request.ApiArguments)
	if err != nil {
		PrintError(fmt.Sprintf("Unable to insert account. %s", err))
		return APIFail(InsertFail)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		PrintError(fmt.Sprintf("Unable to get insertid. %s", err))
		return APIFail(NoInsertID)
	}

	return APISuccess(lastInsertId)
}

func GetCreateAccountStatus(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey", "ID" }); !ok {
		return APIFail(err)
	}

	var req AccountCreationRequest
	qs := `SELECT * FROM AccountCreationRequests JOIN Users ON Users.ID = AccountCreationRequests.UserID WHERE User.ApiKey = :ApiKey AND AccountCreationRequests.ID = :ID`
	if err := SelectOne(&req, qs, request.ApiArguments); err != nil {
		PrintError(fmt.Sprintf("Unable to retrieve account creation request. %s", err))
		return APIFail(SelectFail)
	}

	PrintSuccess(fmt.Sprintf("Successfully retrieved account creation request %+v", req))
	return APISuccess(req)
}

func CreateAccount(request *APIRequest) *APIResponse {
	if ok, err := ValidateRequest(request, []string { "ApiKey", "Email", "DisplayName", "Password", "Age" }); !ok {
		return APIFail(err)
	}

	qs := `INSERT INTO AccountCreationRequests(UserID, Email, DisplayName, Password, Age, StartTime) VALUES((SELECT ID FROM Users WHERE ApiKey = :ApiKey), :Email, :DisplayName, :Password, :Age, :StartTime)`
	request.ApiArguments["StartTime"] = Timestamp()

	result, err := Insert(qs, request.ApiArguments)
	if err != nil {
		PrintError(fmt.Sprintf("Unable to insert account creation request. %s", err))
		return APIFail(InsertFail)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		PrintError(fmt.Sprintf("Unable to get insertid. %s", err))
		return APIFail(NoInsertID)
	}

	return APISuccess(lastInsertId)
}
