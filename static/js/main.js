/*****************************************************************************
 * 
 * Copyright 2020 gtindo.dev
 * 
 * Author: GUEKENG TINDO Yvan
 * 
 * Key Words: Entry point
 * 
 * Purpose: Contains main function which is the entry point of application
 * 
 ****************************************************************************/

import { startBroadcast, stopBroadcast } from './webrtc.js';


function main(){
  let startButton = document.getElementById('startButton');
  startButton.addEventListener('click', startBroadcast);

  let stopButton = document.getElementById('stopButton');
  stopButton.addEventListener("click", stopBroadcast)
}


// run main function
main();
