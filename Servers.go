package main

func GetAllServers(request *APIRequest) *APIResponse {
	return ValidateAndSelectAll(
		request, 
		&[]Server{}, 
		"SELECT Servers.* FROM Servers JOIN Users ON Users.ID = Servers.UserID WHERE Users.ApiKey = :ApiKey", 
		[]string { "ApiKey" }, 
	)
}

func GetServer(request *APIRequest) *APIResponse {
	return ValidateAndSelectOne(
		request,
		&Server{}, 
		"SELECT Servers.* FROM Servers JOIN Users ON Users.ID = Servers.UserID WHERE Users.ApiKey = :ApiKey", 
		[]string { "ApiKey", "ID" }, 
	)
}

func AddServer(request *APIRequest) *APIResponse {
	return ValidateAndInsert(
		request,
		"INSERT INTO Servers(UserID, IPAddress, Name) VALUES((SELECT ID FROM Users WHERE ApiKey = :ApiKey), :IPAddress, :Name)",
		[]string { "ApiKey", "IPAddress", "Name" },
	)
}

func UpdateServer(request *APIRequest) *APIResponse {
	return ValidateAndUpdateFromStruct(
		request,
		&Server{},
		"Servers",
		"UPDATE Servers JOIN Users ON Users.ID = Servers.UserID SET",
		" WHERE Users.ApiKey = :ApiKey AND Servers.ID = :ID",
		[]string { "ApiKey", "ID" },
		"UserID",
	)
}

func DeleteServer(request *APIRequest) *APIResponse {
	return ValidateAndDelete(
		request,
		`DELETE Servers FROM Servers JOIN Users ON Users.ID = Servers.UserID WHERE Users.ApiKey = :ApiKey AND Servers.ID = :ID`,
		[]string { "ApiKey", "ID" },
	)
}

func GetClientsForServer(request *APIRequest) *APIResponse {
	return ValidateAndSelectAll(
		request, 
		&[]Client{}, 
		"SELECT Clients.* FROM Clients JOIN Servers ON Clients.ServerID = Servers.ID JOIN Users ON Servers.UserID = Users.ID WHERE Users.ApiKey = :ApiKey AND Servers.ID = :ID", 
		[]string { "ApiKey", "ID" },
	)
}