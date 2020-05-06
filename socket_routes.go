package main

import (
	"github.com/gorilla/websocket"
)

func GetSocketHandler(route string) func(*websocket.Conn, map[string]string, chan string) {
	switch route {
	case "offer":
		return RTCConnector
	case "candidate":
		return RTCCandidateReceiver
	}

	return nil
}
