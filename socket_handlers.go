package main

import (
	"encoding/json"

	"github.com/pion/webrtc"

	"github.com/gorilla/websocket"
)

func OfferHandler(conn *websocket.Conn, route string, data []byte, c chan string) {
	offerData := SDPData{}
	msgType := 1
	var res []byte

	err := json.Unmarshal(data, &offerData)
	if err != nil {
		res = ErrorResponse(route, ERR400)
		conn.WriteMessage(1, res)
	}

	offer := webrtc.SessionDescription{
		SDP:  offerData.SDP,
		Type: webrtc.SDPType(1), // offer
	}

	// Offer is good
	// Create RTCPeerConnection
	pc, _ := webrtc.NewPeerConnection(webrtc.Configuration{})

	// Allow server to receive audio and video track
	_, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo)
	_, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio)
	if err != nil {
		res = ErrorResponse(route, ERR500)
		conn.WriteMessage(msgType, res)
		panic(err)
	}

	// Set remote description
	err = pc.SetRemoteDescription(offer)
	if err != nil {
		res = ErrorResponse(route, ERR500)
		conn.WriteMessage(msgType, res)
		panic(err)
	}

	// Create answer
	answer, err := pc.CreateAnswer(nil)
	pc.SetLocalDescription(answer)

	res = SuccessResponse(route, SDPData{
		DType: "answer",
		SDP:   answer.SDP,
		UUID:  offerData.UUID,
	})

	// Send answer to user
	err = conn.WriteMessage(msgType, res)
	if err != nil {
		panic(err)
	}

	// Start connector goroutine
	go RTCConnector(conn, pc, offerData, c)
}

func CandidateHandler(conn *websocket.Conn, route string, data []byte, c chan string) {
	cd := CandidateData{}

	err := json.Unmarshal(data, &cd)
	if err != nil {
		res := ErrorResponse(route, ERR400)
		conn.WriteMessage(1, res)
	}

	// c <- cd.ToString()
}
