package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc"
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

func BroadcastHandler(w http.ResponseWriter, r *http.Request) {
	var data SDPData
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
	}

	SdpChan <- data.SDP

	var sdpType webrtc.SDPType = 1

	configuration := webrtc.Configuration{}
	description := webrtc.SessionDescription{
		SDP:  data.SDP,
		Type: sdpType,
	}

	connection, err := webrtc.NewPeerConnection(configuration)
	if err != nil {
		fmt.Println(err)
	}

	err = connection.SetRemoteDescription(description)
	if err != nil {
		fmt.Println(err)
	}

	answer, err := connection.CreateAnswer(nil)

	resData := SDPData{
		DType: "response",
		SDP:   answer.SDP,
		UUID:  data.UUID,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resData)
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

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// parse the message

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			fmt.Fprintf(w, "message response")
		}
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
