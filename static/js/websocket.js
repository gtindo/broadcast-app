/*******************************************************************
 * 
 * Copyright 2020 gtindo.dev
 * 
 * Author: GUEKENG TINDO Yvan
 * 
 * Key Words: websocket, channel, event, listener
 * 
 * Purpose:
 *  Framework on top of websocket to send messages using channels
 * 
 ******************************************************************/


const SEPARTOR = ':---s---:'

export class WS extends EventTarget{

  /**
   * Create websocket object
   * @param {String} url 
   */
  constructor(url){
    super();

    try {
      this.url = new URL(url)
      this.socket = new WebSocket(url);

      let ws = this;
      let socket = this.socket;
      socket.onopen = (event) => {
        console.log("Socket connection established !");

        socket.onmessage = (event) => {
          let data = event.data.split(SEPARTOR);
          let channel = data[0];
          let message = data[1];
          
          let ev = new CustomEvent(channel, {detail : message});
          ws.dispatchEvent(ev);
        }
      }

    } catch(err) {
      throw(err)
    }
  }

  /**
   * Send message to channel
   * 
   * @param {String} channel 
   * @param {String} message 
   */
  send(channel, message) {
    if(typeof(message) !== 'string') throw("Invalid message, must be a string.");
    if(typeof(message) !== 'string') throw("Invalid channel, must be a string");

    let data = `${channel}${SEPARTOR}${message}`;
    this.socket.send(data);
  }
}

/*
ws.socket.onopen = async (event) => {

  console.log("Socket connection opened");

  socket.onmessage = (event) => {
    let res = JSON.parse(event.data);
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
}*/