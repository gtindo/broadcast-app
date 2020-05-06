package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc"
	"github.com/pion/webrtc/pkg/media"
	"github.com/pion/webrtc/pkg/media/ivfwriter"
	"github.com/pion/webrtc/pkg/media/oggwriter"

	"github.com/gorilla/websocket"
)

const (
	rtcpPLIInterval = time.Second * 3
)

func RTCConnector(conn *websocket.Conn, data map[string]string, c chan string) {
	msgType := 1

	req := SDPData{
		DType: data["dtype"],
		SDP:   data["sdp"],
		UUID:  data["uuid"],
	}

	configuration := webrtc.Configuration{}
	offer := webrtc.SessionDescription{
		SDP:  req.SDP,
		Type: webrtc.SDPType(1),
	}

	connection, err := webrtc.NewPeerConnection(configuration)
	if err != nil {
		fmt.Println(err)
	}

	// Allow server to receive 1 video track
	if _, err = connection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		fmt.Println("Codec error !!!")
		fmt.Println(err)
	}

	oggFile, err := oggwriter.New("output.ogg", 48000, 2)
	if err != nil {
		panic(err)
	}
	ivfFile, err := ivfwriter.New("output.ivf")
	if err != nil {
		panic(err)
	}

	// Add Tracke Listener
	connection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := connection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		codec := remoteTrack.Codec()
		if codec.Name == webrtc.Opus {
			fmt.Println("Got Opus track, saving to disk as output.opus (48 kHz, 2 channels)")
			saveToDisk(oggFile, remoteTrack)
		} else if codec.Name == webrtc.VP8 {
			fmt.Println("Got VP8 track, saving to disk as output.ivf")
			saveToDisk(ivfFile, remoteTrack)
		}

		/*
			// Create a local track, all our SFU clients will be fed via this track
			localTrack, newTrackErr := connection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
			if newTrackErr != nil {
				panic(newTrackErr)
			}

			rtpBuf := make([]byte, 1400)
			for {
				i, readErr := remoteTrack.Read(rtpBuf)
				if readErr != nil {
					panic(readErr)
				}
				fmt.Println("Reading...")

				// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
				if _, err = localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
					panic(err)
				}
			}
		*/

	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	connection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateConnected {
			fmt.Println("Ctrl+C the remote client to stop the demo")
		} else if connectionState == webrtc.ICEConnectionStateFailed ||
			connectionState == webrtc.ICEConnectionStateDisconnected {
			closeErr := oggFile.Close()
			if closeErr != nil {
				panic(closeErr)
			}

			closeErr = ivfFile.Close()
			if closeErr != nil {
				panic(closeErr)
			}

			fmt.Println("Done writing media files")
			os.Exit(0)
		}
	})

	err = connection.SetRemoteDescription(offer)
	if err != nil {
		fmt.Println(err)
	}

	answer, err := connection.CreateAnswer(nil)
	connection.SetLocalDescription(answer)

	response := SocketResponse{
		Event: "offer",
		Data: SDPData{
			DType: "answer",
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

	for i := 0; i < 2; i++ {
		candidate := <-c
		s := strings.Split(candidate, "+++")
		sdpMid := s[2]
		intI, _ := strconv.Atoi(s[3])
		sdpIndex := uint16(intI)
		if s[0] == req.UUID {
			iceInit := webrtc.ICECandidateInit{
				Candidate:        s[1],
				SDPMid:           &sdpMid,
				SDPMLineIndex:    &sdpIndex,
				UsernameFragment: s[4],
			}
			err = connection.AddICECandidate(iceInit)
			fmt.Println("Candidate Client : ", iceInit.Candidate)
		}
	}

}

func RTCCandidateReceiver(conn *websocket.Conn, data map[string]string, c chan string) {
	// Send ice candidate to RTCConnector goroutine global CandidateChan channel
	candidate := data["uuid"] + "+++" + data["candidate"] + "+++" + data["sdpMid"] + "+++" + data["sdpMLineIndex"] + "+++" + data["usernameFragment"]
	c <- candidate
}

func saveToDisk(i media.Writer, track *webrtc.Track) {
	defer func() {
		if err := i.Close(); err != nil {
			panic(err)
		}
	}()

	for {
		rtpPacket, err := track.ReadRTP()
		if err != nil {
			panic(err)
		}
		if err := i.WriteRTP(rtpPacket); err != nil {
			panic(err)
		}
	}
}
