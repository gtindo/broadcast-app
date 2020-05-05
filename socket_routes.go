package main

import "github.com/gorilla/websocket"

func GetSocketHandler(route string) func(*websocket.Conn, map[string]string) {
	switch route {
	case "offer":
		return RTCConnector
	}

	return nil
}
