import { startBroadCast, rtcConnection } from './main.js';

let socket = new WebSocket("ws://localhost:4000/socket/");

const startButton = document.querySelector("#startButton");

socket.onopen = async (event) => {

  console.log("Socket connection opened");

  startButton.addEventListener("click", () => {
    startBroadCast(socket);
  })

  socket.onmessage = (event) => {
    let data = JSON.parse(event.data);
    console.log(data);

    //rtcConnection.setRemoteDescription()
  };
}