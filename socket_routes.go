package main

import (
	"github.com/gorilla/websocket"
)

func GetSocketHandler(route string) func(*websocket.Conn, string, []byte, chan string) {
	switch route {
	case "offer":
		return OfferHandler
	case "candidate":
		return CandidateHandler
	}

	return nil
}
