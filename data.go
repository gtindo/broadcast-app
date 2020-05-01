package main

type IndexPageData struct {
	PageTitle string
}

var SdpChan chan string = make(chan string)

type SDPData struct {
	DType string `json:"dtype"`
	SDP   string `json:"sdp"`
	UUID  string `json:"uuid"`
}
