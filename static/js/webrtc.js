/**********************************************************************
 * Copyright 2020 gtindo.dev
 * 
 * Author: GUEKENG TINDO Yvan
 * 
 * Key Words: webrtc, media, streaming
 * 
 * Purpose: 
 *  Functions defined in this file are used to establish a
 *  RTCPeerConnection with a remote host. We follow these steps :
 *   1. Get media stream from browser
 *   2. Create RTCPeerConnection object
 *   3. Add media track to connection
 *   4. Create an offer
 *   5. Set local session description
 *   6. Send this description to remote using signaling service
 *   7. Get answer from signaling service
 *   8. Set remote description
 *   9. Send ICE candidates to remote with signaling service
 *   10. Receive remote ICE candidates
 *   11. Add remote candidates to RTCPeerConnection
 *   12. Listening to changes on connection status
 * 
 *********************************************************************/

import { randomString, setStopBehavior, startTimer } from './utils.js';
import { WS } from './websocket.js';


// Establish websocket connection
const WS_URL = ("ws://localhost:4000/socket/");
let ws = new WS(WS_URL);
ws.addEventListener("error", (err) => console.log(err))
ws.addEventListener("download_file", (event) => {
  let res = JSON.parse(event.detail);
  let message = document.getElementById("message");
  message.innerHTML = `<a href="${res.data.message}">Your file</a>`
});


// unique id of client
export const CLIENT_UUID = randomString(25);

// Peer Connection
const WEBRTC_CONFIGURATION = null;
let PC = new RTCPeerConnection(WEBRTC_CONFIGURATION);

// Set up to exchange audio and video
const OFFER_OPTIONS = {
  offerToReceiveVideo: 1,
  offerToREceiveAudio: 1
};

// Error message
const CONNECTION_ERROR_MSG = "Error while establishing connection with stream server."

// Set constraints on stream
const MEDIA_CONSTRAINTS = {
  video: {
    aspectRatio: { exact: 16/9 },
    facingMode: { ideal: "user" },
    frameRate: { ideal: 30 },
    height: { ideal: 480 },
    width: { ideal: 720 },
    resizeMode: { ideal: "crop-and-scale"}
  },
  audio: true
};

// Video element where stream will be placed
const localVideo = document.querySelector("#local-video");

// Layer before video element
let vidBefore = document.getElementById('vid-before');

// Local stream where stream will be placed
let localStream;

// Buttons
let startButton = document.getElementById('startButton');
let stopButton = document.getElementById('stopButton');

// Stream call backs

/**
 * Handle success by adding the mediastream to video element
 * @param {MediaStream} mediaStream 
 */
function gotLocalMediaStream(mediaStream)
{ 
  localStream = mediaStream;
  localVideo.srcObject = mediaStream; 

  // Set gray layer on video element
  vidBefore.className = "vid-before";
}


/**
 * Handle error by logging a message to console with error message
 * @param {MediaStreamError} error 
 */
function handleLocalMediaStreamError(error) {
  console.log('navigator.getUserMedia error: ', error);
  alert("Can't get your camera or microphone, please accept permissions in order to use the app.");
}


/**
 * Get usermedia and initialize a Peer Connection with streaming server.
 * 
 * If an error occurs while getting user media it will show an error to user and ask him
 * to accept permissions.
 * 
 * The video will be blurred until the connection has been established.
 */
export async function startBroadcast() {
  try{
    // Initializes media stream
    let mediaStream = await navigator.mediaDevices.getUserMedia(MEDIA_CONSTRAINTS);
    gotLocalMediaStream(mediaStream);
    await rtcConnect();
  } catch (err) {
    handleLocalMediaStreamError(err);
  }
}


/**
 * Log session description error and show error message to user
 * @param {Error} err 
 */
function setSessionDescriptionError(err){
  
}


/**
 * Listens to connection status :
 * 
 * @param {Event} event 
 */
function handlePCStateChange(event){
  console.log(PC.iceConnectionState);
  if(PC.iceConnectionState == "connected"){
    // Remove gray layer
    vidBefore.className = "vid-before-start";
    
    // Remove start button on screen
    startButton.style.display = "none";

    // Display Stop button
    stopButton.style.display = "inline";

    // Start timer
    startTimer();
  }
}


/**
 * Listens to changes on ice candidate 
 * and when a new candidate is discovered it is send to remote.
 * 
 * @param {Event} event 
 */
function handleIceCandidate(event){
  let candidate = event.candidate
  if(candidate !== null) {
    console.log(candidate)
    let data = {
      "event": "candidate",
      "data": {
        "uuid": CLIENT_UUID,
        "candidate": candidate.candidate,
        "sdpMid": candidate.sdpMid,
        "sdpMLineIndex": ""+candidate.sdpMLineIndex,
        "usernameFragment": candidate.usernameFragment
      }
    }
    // send candidate to remote
  }
}


/**
 * 
 * @param {Event} data 
 */
function handleOfferResponse(event){
  let detail = JSON.parse(event.detail)

  let sessionInit = {
    type: detail.data.dtype,
    sdp: detail.data.sdp
  }
  console.log(sessionInit);
  PC.setRemoteDescription(new RTCSessionDescription(sessionInit));
}

// Create peer connection and behavior
/**
 * Establish connection with Streaming server
 * - Add track to peerc connection
 * - set local description
 * - send local description to server
 * - set connection state listener
 * - set ICE candidate Listener
 * 
 * @param {WebSocket} socket 
 */
export async function rtcConnect(){
  // Write display connecting on start button
  startButton.innerHTML = "Connecting..."

  // Add localstream to connection and create offer to connect
  let tracks = localStream.getTracks()
  for(let track of tracks){
    PC.addTrack(track)
  }

  try {
    let description = await PC.createOffer(OFFER_OPTIONS);
    PC.setLocalDescription(description);
    
    let data = {
        "dtype": "offer",
        "sdp": description.sdp,
        "uuid": CLIENT_UUID
    }

    ws.send("offer", JSON.stringify(data));

    PC.onconnectionstatechange = handlePCStateChange;

    ws.addEventListener("offer", handleOfferResponse);

  } catch (err) {
    console.log(CONNECTION_ERROR_MSG)
    console.log(err)

    alert(CONNECTION_ERROR_MSG);

    // remove Gray layer
    setStopBehavior();
  }
}


export function stopBroadcast(){
  PC.close();
  PC = new RTCPeerConnection(WEBRTC_CONFIGURATION);

  setStopBehavior();
}