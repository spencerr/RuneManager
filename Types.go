package main

import (
	"github.com/labstack/echo"
)

const (
	InsertFail = "insert-fail"
	SelectFail = "select-fail"
	UpdateFail = "update-fail"
	DeleteFail = "delete-fail"
	NoInsertID = "no-insert-id"
	NoRowsAffected = "no-rows-affected"
)

var (
	Functions = map[string] func(*APIRequest) *APIResponse {
		"no-route": NoRoute,
		"get-all-accounts": GetAccounts,
		"add-account": AddAccount,

		"get-account": GetAccount,
		"update-account": UpdateAccount,
		"delete-account": DeleteAccount,
		"lock-account": LockAccount,
		"ban-account": BanAccount,

		"create-account": CreateAccount,
		"get-create-account-status": GetCreateAccountStatus,
		"reset-password": ResetAccountPassword,
		"get-reset-password-status": GetResetPasswordStatus,
	}
)

type Account struct {
	ID			int64	`json:"id" db:"ID" static:"true"`
	UserID		int64	`json:"userid" db:"UserID" static:"true"`
	Email		string	`json:"email" db:"Email"`
	Password	string	`json:"password" db:"Password"`
	Locks		int64	`json:"locks" db:"Locks"`
	Locked		bool	`json:"locked" db:"Locked"`
	Banned		bool	`json:"banned" db:"Banned"`
}

type User struct {
	ID			int64		`json:"id" db:"ID" static:"true"`
	Email		string		`json:"email" db:"Email" static:"true"`
	Password	string		`json:"password" db:"Password"`
	ApiKey		string		`json:"email" db:"ApiKey"`
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
	ID			int64	`json:id db:"ID" static:"true"`
	UserID		int64	`json:userid db:"UserID" static:"true"`
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
	ID			int64 	`json:"id" db:"ID" static:"true"`
	UserID		int64	`json:userid db:"UserID" static:"true"`
	AccountID	int64	`json:accountid db:"AccountID" static:"true"`
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

type Server struct {
	ID			int64 	`json:"id" db:"ID" static:"true"`
	UserID		int64	`json:userid db:"UserID" static:"true"`
	IPAddress	int64	`json:ip_address db:"IPAddress" static:"true"`
	Name		string	`json:name`
}

func BindServer(c echo.Context) *Server {
	server := new(Server)
	if err := c.Bind(server); err != nil {
		return nil
	}

	if id, ok := ParamToInt64(c, "id"); ok {
		server.ID = id
	}

	return server
}

type Client struct {
	ID					int64 	`json:"id" db:"ID" static:"true"`
	ServerID			int64	`json:"serverid" db:"ServerID" static:"true"`
	AccountID			int64	`json:"accountid" db:"AccountID"`
	ScriptName			string	`json:"script_name" db:"ScriptName"`
	ScriptArguments		string	`json:"script_arguments" db:"ScriptArguments"`
	World				int64	`json:"world" db:"World"`
}

func BindClient(c echo.Context) *Client {
	client := new(Client)
	if err := c.Bind(client); err != nil {
		return nil
	}

	if id, ok := ParamToInt64(c, "id"); ok {
		client.ID = id
	}

	return client
}