package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Home page")
}

func IndexHandler(res http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	data := IndexPageData{
		PageTitle: "Broadcast App | Home",
	}

	tmpl.Execute(res, data)
}

// websocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var data SocketMsg
		var response SocketResponse
		var errMsg ErrorMsg

		err = json.Unmarshal(msg, &data)
		if err != nil {
			errMsg = ErrorMsg{Message: "Bad message", Code: "ERR400"}
			response = SocketResponse{
				Event:  "",
				Data:   nil,
				Error:  errMsg,
				Status: false,
			}

			// TODO : Log errors in file
			res, _ := json.Marshal(response)
			conn.WriteMessage(msgType, res)
		}

		route := GetSocketHandler(data.Event)
		route(conn, data.Data)
	}
}

/*
1. Server receive sdp through websocket connection
2. Send sdp on connection channel
3. Connection goroutine receive sdp
4. Create RTC connection
5. send answer to websocket handler
6. websocket send response to client
*/
