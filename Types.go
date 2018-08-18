package main

import (
	"github.com/labstack/echo"
)


type Account struct {
	ID			int64	`json:id`
	UserID		int64	`json:userid`
	Email		string	`json:email`
	Password	string	`json:password`
	Locks		int64	`json:locks`
	Banned		bool	`json:banned`
}

func BindAccount(c echo.Context) *Account {
	account := new(Account)
	if err := c.Bind(account); err != nil {
		return nil
	}

	if id, ok := ParamToInt64(c, "id"); ok {
		account.ID = id
	}

	return account
}

type AccountCreationRequest struct {
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

func BindAccountCreationRequest(c echo.Context) *AccountCreationRequest {
	request := new(AccountCreationRequest)
	if err := c.Bind(request); err != nil {
		return nil
	}

	if id, ok := ParamToInt64(c, "id"); ok {
		request.ID = id
	}

	return request
}

type PasswordResetRequest struct {
	ID			int64 	`json:"id" form:"id" query:"id"`
	UserID		int64	`json:userid`
	AccountID	int64	`json:accountid`
	NewPassword	string	`json:new_password`
	RequestURL	string	`json:request_url`
	StartTime	int64	`json:start_time`
	EndTime		int64	`json:end_time`
	Status		string	`json:status`
}

func BindPasswordResetRequest(c echo.Context) *PasswordResetRequest {
	request := new(PasswordResetRequest)
	if err := c.Bind(request); err != nil {
		return nil
	}

	if id, ok := ParamToInt64(c, "id"); ok {
		request.ID = id
	}

	return request
}

type RequestResult struct {
	Success		bool	`json: success`
	Message		string	`json: message`
	Result 		interface{} `json: result`
}