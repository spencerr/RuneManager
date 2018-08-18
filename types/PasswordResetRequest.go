package types

type PasswordResetRequest struct {
	ID			int64 	`json:id`
	UserID		int64	`json:userid`
	AccountID	int64	`json:accountid`
	NewPassword	string	`json:new_password`
	RequestURL	string	`json:request_url`
	StartTime	int64	`json:start_time`
	EndTime		int64	`json:end_time`
	Status		string	`json:status`
}

func New() PasswordResetRequest {
	return &PasswordResetRequest
}

func Bind(c echo.Context) PasswordResetRequest {
	request := new(PasswordResetRequest)
	if err = c.Bind(request); err != nil {
		return nil
	}

	return request
}