// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"fmt"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients 		map[*Socket]bool
	send 			chan *Message
	register 		chan *Socket
	unregister 		chan *Socket
	handleMessage 	chan *SocketMessage
}

type Message struct {
	clientid 	string
	text 		string
}


type ClientResponse struct {
	clients 	[]ClientInformation
}

type ClientInformation struct {
	id					int64
	accountid 			int64
	account_email		string
	account_password	string
	script_name			string
	script_arguments 	string
}

type SocketMessage struct {
	Socket		*Socket
	Message		[]byte
}

type APIRequest struct {
	ApiRoute		string
	ApiArguments	map[string]interface{}
}

type APIResponse struct {
	Success			bool
	Result			interface{}
}

var hub = newHub()

func newHub() *Hub {
	return &Hub{
		send:			make(chan *Message),
		register:   	make(chan *Socket),
		unregister: 	make(chan *Socket),
		clients:    	make(map[*Socket]bool),
		handleMessage: 	make(chan *SocketMessage),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <- h.send:
			h.sendTo(message.clientid, message.text)

		case data := <- h.handleMessage:
			request := gjson.GetBytes(data.Message, "Type").String()
			if request == "APIRequest" {
				var apiRequest APIRequest
				if json.Unmarshal(data.Message, &apiRequest) == nil {
					fmt.Printf("Attempting to request with data %+v", apiRequest)
					if bytes, err := json.Marshal(Functions[apiRequest.ApiRoute](&apiRequest)); err == nil {
						data.Socket.send <- bytes
					}
				}
			}
		}
	}
}

func (h *Hub) sendTo(clientid string, message string) {
	for client := range h.clients {
		if client.id == clientid {
			select {
			case client.send <- []byte(message):
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}