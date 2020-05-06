import { startBroadCast, rtcConnection, clientUUID } from './main.js';

let socket = new WebSocket("ws://localhost:4000/socket/");

const startButton = document.querySelector("#startButton");

socket.onopen = async (event) => {

  console.log("Socket connection opened");

  startButton.addEventListener("click", () => {
    startBroadCast(socket);
  })

  socket.onmessage = (event) => {
    let res = JSON.parse(event.data);
    console.log(res);

    switch(res.event){
      case "offer":
        let sessionInit = {
          type: res.data.dtype,
          sdp: res.data.sdp
        }
        rtcConnection.setRemoteDescription(new RTCSessionDescription(sessionInit), function() {
          console.log(rtcConnection.remoteDescription);
        }, (err) => console.log(err));
        break;

      case "candidate":
        break;
    }
    

  };


  rtcConnection.oniceconnectionstatechange = e => console.log(rtcConnection.iceConnectionState);
  
  rtcConnection.onicecandidate = event => {
    if(event.candidate !== null){
      console.log(event.candidate)
      let data = {
        "event": "candidate",
        "data": {
          "uuid": clientUUID,
          "candidate": event.candidate.candidate,
          "sdpMid": event.candidate.sdpMid,
          "sdpMLineIndex": ""+event.candidate.sdpMLineIndex,
          "usernameFragment": event.candidate.usernameFragment
        }
      }
      socket.send(JSON.stringify(data))
      console.log("Candidate send")
      console.log(rtcConnection.iceConnectionState);
    }
  }
}