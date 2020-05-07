package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// Home page
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

// Websocket server
func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	channel := make(chan string)

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		msgtab := strings.Split(string(msg), SEPARATOR)

		if len(msgtab) != 2 {
			conn.WriteMessage(msgType, []byte("error"+SEPARATOR+"WS class error, bad message sent."))
			panic("Invalid message received")
		}

		route := msgtab[0]
		data := []byte(msgtab[1])

		handler := GetSocketHandler(route)
		handler(conn, route, data, channel)
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
