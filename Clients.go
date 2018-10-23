package main

func GetAllClients(request *APIRequest) *APIResponse {
	return ValidateAndSelectAll(
		request, 
		&[]Client{}, 
		"SELECT Clients.* FROM Clients JOIN Users ON Users.ID = Clients.UserID WHERE Users.ApiKey = :ApiKey", 
		[]string { "ApiKey" }, 
	)
}

func GetClient(request *APIRequest) *APIResponse {
	return ValidateAndSelectOne(
		request,
		&Client{}, 
		"SELECT Clients.* FROM Clients JOIN Users ON Users.ID = Clients.UserID WHERE Users.ApiKey = :ApiKey AND Clients.ID = :ID", 
		[]string { "ApiKey", "ID" }, 
	)
}

func AddClient(request *APIRequest) *APIResponse {
	return ValidateAndInsert(
		request,
		"INSERT INTO Clients(UserID, ServerID, AccountID, ScriptName, ScriptArguments, World) VALUES(:UserID, (SELECT Accounts.ID FROM Accounts WHERE Accounts.UserID = :UserID AND Accounts.ID = :AccountID), (SELECT Servers.ID FROM Servers WHERE Servers.UserID = :UserID AND Servers.ID = :ID), :ScriptName, :ScriptArguments, :World)",
		[]string { "ApiKey", "IPAddress", "Name" },
	)
}

func UpdateClient(request *APIRequest) *APIResponse {
	return ValidateAndUpdateFromStruct(
		request,
		&Server{},
		"Servers",
		"UPDATE Clients JOIN Users ON Users.ID = Clients.UserID SET",
		" WHERE Users.ApiKey = :ApiKey AND Clients.ID = :ID",
		[]string { "ApiKey", "ID" },
		"UserID",
	)
}

func DeleteClient(request *APIRequest) *APIResponse {
	return ValidateAndDelete(
		request,
		`DELETE Clients FROM Clients JOIN Users ON Users.ID = Clients.UserID WHERE Users.ApiKey = :ApiKey AND Clients.ID = :ID`,
		[]string { "ApiKey", "ID" },
	)
}