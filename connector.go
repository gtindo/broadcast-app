package main

import (
	"encoding/json"
	"fmt"

	"github.com/pion/webrtc"

	"github.com/gorilla/websocket"
)

func RTCConnector(conn *websocket.Conn, data map[string]string) {
	msgType := 1

	req := SDPData{
		DType: data["dtype"],
		SDP:   data["sdp"],
		UUID:  data["uuid"],
	}

	configuration := webrtc.Configuration{}
	description := webrtc.SessionDescription{
		SDP:  req.SDP,
		Type: webrtc.SDPType(1),
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

	response := SocketResponse{
		Event: "offer",
		Data: SDPData{
			DType: "response",
			SDP:   answer.SDP,
			UUID:  req.UUID,
		},
		Error:  ErrorMsg{},
		Status: true,
	}

	res, _ := json.Marshal(response)

	if err := conn.WriteMessage(msgType, res); err != nil {
		fmt.Println(err)
	}
}
