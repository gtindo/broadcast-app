import { startBroadCast } from './main.js'

let socket = new WebSocket("ws://localhost:4000/socket/");

const startButton = document.querySelector("#startButton");

socket.onopen = async (event) => {

  console.log("Socket connection opened");

  startButton.addEventListener("click", () => {
    socket.send("Ceci est un message");
  })

  socket.onmessage = (event) => {
    console.log(event.data);
  };
}