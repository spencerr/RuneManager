package main

func GetAccounts(request *APIRequest) *APIResponse {
	return ValidateAndSelectAll(
		request,
		&[]Account{},
		"SELECT Accounts.* FROM Accounts JOIN Users ON Users.ID = Accounts.UserID WHERE Users.ApiKey = :ApiKey",
		[]string { "ApiKey" },
	)
}

func AddAccount(request *APIRequest) *APIResponse {
	return ValidateAndInsert(
		request,
		"INSERT INTO Accounts(UserID, Email, Password) VALUES((SELECT ID FROM Users WHERE ApiKey = :ApiKey), :Email, :Password)",
		[]string { "ApiKey", "Email", "Password" },
	)
}

func GetAccount(request *APIRequest) *APIResponse {
	return ValidateAndSelectOne(
		request,
		&Account{},
		"SELECT Accounts.* FROM Accounts JOIN Users ON Users.ID = Accounts.UserID WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID",
		[]string { "ApiKey", "ID" },
	)
}

func UpdateAccount(request *APIRequest) *APIResponse {
	return ValidateAndUpdateFromStruct(
		request,
		&Account{},
		"Accounts",
		"UPDATE Accounts JOIN Users ON Users.ID = Accounts.UserID SET",
		" WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID",
		[]string { "ApiKey", "ID" },
		"UserID",
	)
}

func DeleteAccount(request *APIRequest) *APIResponse {
	return ValidateAndDelete(
		request,
		`DELETE Accounts FROM Accounts JOIN Users ON Users.ID = Accounts.UserID WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID`,
		[]string { "ApiKey", "ID" },
	)
}

func LockAccount(request *APIRequest) *APIResponse {
	return ValidateAndUpdate(
		request,
		`UPDATE Accounts JOIN Users ON Users.ID = Accounts.UserID SET Locks = Locks + 1 WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID`,
		[]string { "ApiKey", "ID" },
	)
}

func BanAccount(request *APIRequest) *APIResponse {
	return ValidateAndUpdate(
		request,
		`UPDATE Accounts JOIN Users ON Users.ID = Accounts.UserID SET Banned = TRUE WHERE Users.ApiKey = :ApiKey AND Accounts.ID = :ID`,
		[]string { "ApiKey", "ID" },
	)
}

func GetResetPasswordStatus(request *APIRequest) *APIResponse {
	return ValidateAndSelectOne(
		request,
		&PasswordResetRequest{},
		`SELECT PasswordResetRequests.* FROM PasswordResetRequests JOIN Users ON Users.ID = PasswordResetRequests.UserID WHERE User.ApiKey = :ApiKey AND PasswordResetRequest.ID = :ID`,
		[]string { "ApiKey", "ID" },
	)
}

func ResetAccountPassword(request *APIRequest) *APIResponse {
	request.ApiArguments["StartTime"] = Timestamp()

	return ValidateAndInsert(
		request,
		`INSERT INTO PasswordResetRequests(UserID, AccountID, NewPassword, StartTime) VALUES((SELECT ID FROM Users WHERE ApiKey = :ApiKey), :ID, :NewPassword, :StartTime)`,
		[]string { "ApiKey", "NewPassword", "ID" },
	)
}

func GetCreateAccountStatus(request *APIRequest) *APIResponse {
	return ValidateAndSelectOne(
		request,
		&AccountCreationRequest{},
		`SELECT AccountCreationRequests.* FROM AccountCreationRequests JOIN Users ON Users.ID = AccountCreationRequests.UserID WHERE User.ApiKey = :ApiKey AND AccountCreationRequests.ID = :ID`,
		[]string { "ApiKey", "ID" },
	)
}

func CreateAccount(request *APIRequest) *APIResponse {
	return ValidateAndInsert(
		request,
		`INSERT INTO AccountCreationRequests(UserID, Email, DisplayName, Password, Age, StartTime) VALUES((SELECT ID FROM Users WHERE ApiKey = :ApiKey), :Email, :DisplayName, :Password, :Age, :StartTime)`,
		[]string { "ApiKey", "Email", "DisplayName", "Password", "Age" },
	)
}
