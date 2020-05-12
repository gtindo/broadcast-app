package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc"
	"github.com/pion/webrtc/pkg/media/ivfwriter"
	"github.com/pion/webrtc/pkg/media/oggwriter"
)

const (
	rtcpPLIInterval = time.Second * 3
)

func RTCConnector(ws *websocket.Conn, pc *webrtc.PeerConnection, offerData SDPData, c chan string) {

	basePath := "static/audio/"
	video := basePath + offerData.UUID + ".ivf"
	audio := basePath + offerData.UUID + ".ogg"

	oggFile, err := oggwriter.New((audio), 48000, 2)
	if err != nil {
		panic(err)
	}

	ivfFile, err := ivfwriter.New(video)
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

		writingStatus := 0

		if connectionState == webrtc.ICEConnectionStateConnected {
			fmt.Println("Stream started for - ", offerData.UUID)
			writingStatus = 1
		} else if connectionState == webrtc.ICEConnectionStateFailed {
			fmt.Println("Peer connection failed to reconnect.")
		} else if connectionState == webrtc.ICEConnectionStateDisconnected {
			fmt.Printf("Peer connection State has changed %s \n", connectionState.String())
			closeErr := oggFile.Close()
			if closeErr != nil {
				panic(closeErr)
			}

			closeErr = ivfFile.Close()
			if closeErr != nil {
				panic(closeErr)
			}

			fmt.Println("Done writing media files for -", offerData.UUID)
			writingStatus += 1
		}

		if writingStatus == 2 {
			pre_video := basePath + "pre_" + offerData.UUID + ".ivf"
			output := basePath + offerData.UUID + ".webm"
			cmd1 := exec.Command("ffmpeg", "-i", video, "-filter:v", "setpts=2*PTS", pre_video)
			cmd2 := exec.Command("ffmpeg", "-i", pre_video, "-i", audio, "-c", "copy", output)

			data := make(map[string]string)
			data["message"] = "Generating your video file...."
			res := SuccessResponse("download_start", data)

			ws.WriteMessage(1, res)
			fmt.Println("Generating video file for ", offerData.UUID)
			err = cmd1.Run()
			fmt.Println(err)
			err = cmd2.Run()
			fmt.Println(err)
			fmt.Println("File generation Finished - ", offerData.UUID)

			data["message"] = "/" + output
			res = SuccessResponse("download_file", data)
			ws.WriteMessage(1, res)
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
