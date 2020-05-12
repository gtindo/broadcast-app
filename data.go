package main

import (
	"errors"
	"strings"
)

const ERR400 = "ERR400"
const ERR500 = "ERR500"

var ERRORS_MSG = map[string]string{
	ERR400: "Bad message sent",
	ERR500: "Internal error",
}

const SEPARATOR = ":---s---:"

type IndexPageData struct {
	PageTitle string
}

type SDPData struct {
	DType string `json:"dtype"`
	SDP   string `json:"sdp"`
	UUID  string `json:"uuid"`
}

type CandidateData struct {
	Candidate        string `json:"candidate"`
	SdpMid           string `json:"sdpMid"`
	SdpMLineIndex    string `json:"sdpMLineIndex"`
	UsernameFragment string `json:"usernameFragment"`
	UUID             string `json:"uuid"`
}

func (cd *CandidateData) ToString() string {
	sp := SEPARATOR
	return cd.UUID + sp + cd.Candidate + sp + cd.SdpMid + sp + cd.SdpMLineIndex + sp + cd.UsernameFragment
}

func NewCandidateData(candidate string) (CandidateData, error) {
	sp := SEPARATOR
	values := strings.Split(candidate, sp)
	if len(values) < 5 {
		return CandidateData{}, errors.New("Invalid candidate string.")
	}
	return CandidateData{
		UUID:             values[0],
		Candidate:        values[1],
		SdpMid:           values[2],
		SdpMLineIndex:    values[3],
		UsernameFragment: values[4],
	}, nil
}

type SocketMsg struct {
	Event string            `json:"event"`
	Data  map[string]string `json:"data"`
}

type ErrorMsg struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type SocketResponse struct {
	Data   interface{} `json:"data"`
	Error  ErrorMsg    `json:"error"`
	Status bool        `json:"status"`
}

type Config struct {
	APP_NAME      string `json:"APP_NAME"`
	APP_VERSION   string `json:"APP_VERSION"`
	HTTP_PORT     string `json:"HTTP_PORT"`
	ENV           string `json:"ENV"`
	SSL_KEY_PATH  string `json:"SSL_KEY_PATH"`
	SSL_CERT_PATH string `json:"SSL_CERT_PATH"`
}
