package main

type IndexPageData struct {
	PageTitle string
}

type SDPData struct {
	DType string `json:"dtype"`
	SDP   string `json:"sdp"`
	UUID  string `json:"uuid"`
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
	Event  string      `json:"event"`
	Data   interface{} `json:"data"`
	Error  ErrorMsg    `json:"error"`
	Status bool        `json:"status"`
}

type CadidateData struct {
	Candidate   string      `json:"candidate"`
	Description interface{} `json:"description"`
}
