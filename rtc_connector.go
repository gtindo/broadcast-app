package main

import (
	"fmt"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc"
	"github.com/pion/webrtc/pkg/media/ivfwriter"
	"github.com/pion/webrtc/pkg/media/oggwriter"
)

const (
	rtcpPLIInterval = time.Second * 3
)

func RTCConnector(pc *webrtc.PeerConnection, offerData SDPData, c chan string) {

	basePath := "static/audio/"

	oggFile, err := oggwriter.New((basePath + offerData.UUID + ".ogg"), 48000, 2)
	if err != nil {
		panic(err)
	}

	ivfFile, err := ivfwriter.New(basePath + offerData.UUID + ".webm")
	if err != nil {
		panic(err)
	}

	// Add Track Listener
	pc.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
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

	})

	// Set the handler for ICE pc state
	// This will notify you when the peer has connected/disconnected
	pc.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Peer connection State has changed %s \n", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateConnected {
			fmt.Println("Stream started for - ", offerData.UUID)
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

			fmt.Println("Done writing media files for -", offerData.UUID)
		}
	})
}

/*
	// Create a local track, all our SFU clients will be fed via this track
	localTrack, newTrackErr := pc.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
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
