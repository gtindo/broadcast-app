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


export function getClientID(){
  const id_name = "client_id"
  let id = localStorage.getItem(id_name)

  if(id == undefined){
    id = randomString(25);
    localStorage.setItem(id_name, id);
  }

  return id;
}
