// unique id of client
const clientUUID = "dsaffe2wo0ifsdafowqdsfasf";

// Set constraints on stream

const mediaConstraints = {
  video: {
    optional: [
      { maxWidth: 720 },
      { maxHeight: 480 }
    ]
  },
};

// Video element where stream will be placed
const localVideo = document.querySelector("#local-video");

// Local stream where stream will be placed
let localStream;

// Stream call backs

// Handle success by adding the mediastream to video element
function gotLocalMediaStream(mediaStream)
{ 
  localStream = mediaStream;
  localVideo.srcObject = mediaStream; 
}

// Handle error by logging a message to console with error message
function handleLocalMediaStreamError(error) {
  console.log('navigator.getUserMedia error: ', error)
}

// Initializes media stream
navigator.mediaDevices.getUserMedia(mediaConstraints)
  .then(gotLocalMediaStream)
  .catch(handleLocalMediaStreamError)



/*****************************************************/
const configuration = null;
export let rtcConnection = new RTCPeerConnection(configuration);

// Set up to exchange only video
const offerOptions = {
  offerToReceiveVideo: 1,
};

// Log offer creation and sets peer connection session description
function createdOffer(description){
  return new Promise((resolve, reject) => {
    rtcConnection.setLocalDescription(description)
    .then(() => {
      resolve(description)
      console.log("Local description has been set.")
    }).catch((err) => reject(err));
  });
}

function setSessionDescriptionError(err){
  console.log(err)
  console.log("Error while setting session description")
}


// Create peer connection and behavior
/**
 * 
 * @param {WebSocket} socket 
 */
export function startBroadCast(socket){
  // Add localstream to connection and create offer to connect
  rtcConnection.addTrack(localStream.getTracks()[0])

  rtcConnection.createOffer(offerOptions)
    .then(createdOffer)
    .then((description) => {
      // Send offer to signaling server
      let data = {
        "event": "offer",
        "data": {
          "dtype": "offer",
          "sdp": description.sdp,
          "uuid": clientUUID
        }
      }
      
      socket.send(JSON.stringify(data));
      console.log("Offer send to signaling server.")
    })
    .catch(setSessionDescriptionError);
}


function stopBroadcast(){
  rtcConnection.close();
  rtcConnection = null;
}

/* Protocols
SDP : Session description protocol
ICE : Interactive connectivity establishment
STUN :
TURN :
RTP :
*/

/* Emitter : Web Client
1. Create RTCPeerConnection
2. call RTCPeerConnection.createOffer()
3. call RTCPeerConnection.setLocalDescription()
4. Generate ice candidates with STUN Server or TURN Server (set in servers configuration)
5. Send offer to intended receiver () // Send sdp

12. client receive Answer
13. client call RTCPeerConnection.setRemoteDescription()

END;
*/

/* Receiver : Web Server

6. Create RTCPeerConnection
7. call RTCPeerConnection.setRemoteDescription()
8. call RTCPeerConnection.addTrack()
9. call RTCPeerConnection.createAnswer()
10. call RTCPeerConnection.setLocalDescription()
11. Send answer to the caller

*/