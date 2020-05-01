package main

import "fmt"

func RTCConnector(c chan string) {
	for {
		sdp := <-c
		go NewRTCPeerConnection(sdp)
	}
}

func NewRTCPeerConnection(sdp string) {
	fmt.Println(sdp)
}
