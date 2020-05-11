/*******************************************************************
 * 
 * Copyright 2020 gtindo.dev
 * 
 * Author: GUEKENG TINDO Yvan
 * 
 * Key Words: helpers, utils
 * 
 * Purpose : Contains functions used to made some repetitive tasks
 *
 *******************************************************************/


/**
 * Generate a string with random characters.
 * 
 * @param {Number} nbChars
 * @returns {String} 
 */
export function randomString(nbChars){
  let rndStr = "";
  let chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_";
  let n = chars.length;
  for(let i = 0; i < nbChars; i++){
    let rndIndex = parseInt(Math.random() * nbChars * 100) % n;
    rndStr += chars[rndIndex];
  }

  return rndStr;
}


/**
 * Get time in seconds and rewrite it as a string like :
 * Nh : Nmn : Ns 
 * @param {Number} time
 * @returns {String} 
 */
function formatTimer(time){
  let seconds = time % 60 ;
  
  // In minutes
  let remaining = Math.floor((time - seconds) / 60);

  let minutes = remaining % 60;
  
  let hours = parseInt(remaining / 60)

  return hours + "h : " + minutes + "mn : "+ seconds + "s" 
}

// timer interval
let timerInterval ;

export function startTimer() {
  let timer = document.getElementById("timer");
  let timerBox = document.getElementById("timer-box");
  let time = 0;

  // display timer
  timerBox.className = "timer-box-start"
  
  timerInterval = setInterval(() => {
    time += 1;
    timer.innerHTML = formatTimer(time);
  }, 1000);
}

export function setStopBehavior() {
  let startButton = document.getElementById('startButton');
  let stopButton = document.getElementById('stopButton');
  let vidBefore = document.getElementById('vid-before');
  let timerBox = document.getElementById("timer-box");
  let message = document.getElementById("message")


  vidBefore.className = "vid-before-start";
  startButton.innerHTML = "Broadcast Yourself"
  startButton.style.display = "inline"
  stopButton.style.display = "none";
  clearInterval(timerInterval);
  timerBox.className = "timer-box-stop";
  message.innerHTML = "Generating your video file....";
}