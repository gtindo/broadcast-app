package main

import (
	"encoding/json"

	"github.com/pion/webrtc"
	"github.com/pion/webrtc/pkg/media"
)

// Format error response message
func ErrorResponse(route string, err string) []byte {
	resI := SocketResponse{
		Data: nil,
		Error: ErrorMsg{
			Message: ERRORS_MSG[err],
			Code:    err,
		},
		Status: false,
	}
	resD, _ := json.Marshal(resI)
	pref := route + SEPARATOR

	return []byte(pref + string(resD))
}

// Format success response message
func SuccessResponse(route string, data interface{}) []byte {
	resI := SocketResponse{
		Data:   data,
		Error:  ErrorMsg{},
		Status: true,
	}
	resD, _ := json.Marshal(resI)
	pref := route + SEPARATOR

	return []byte(pref + string(resD))
}

// Save track on disk
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

// ffmpeg -i roqfgrewd_GbvLDaWIDxWOiKx.ivf -filter:v "setpts=2*PTS" output.ivf

// ffmpeg -i roqfgrewd_GbvLDaWIDxWOiKx.ivf -filter:v "setpts=2*PTS" output.ivf && ffmpeg -i output.ivf -i roqfgrewd_GbvLDaWIDxWOiKx.ogg -c copy output.webm
