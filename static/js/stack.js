function isCorrect(expr){
  const n = expr.length;
  const numbers = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "0"];
  const beforeOperator = [...numbers, ")"]
  const afterOperator = [...numbers, "("]
  const operators = ["+", "-", "/", "*", "."];
  const authorisedsChars = [...numbers, ...operators, "(", ")"];
  let leftParenthesis = [];

  for(let i = 0; i < n; i++){
    let c = expr[i];
    // Check if character is authorized
    if(!authorisedsChars.includes(c)){
      console.log(`Character "${c}" is invalid`);
      return false;
    }

    // Check if operators are correct
    if(operators.includes(c)){
      if(i === n-1) {
        console.log(`Character ${c} can't be the last character.`)
        return false;
      }else if(i === 0) {
        console.log(`Character ${C} can't be the first character.`)
        return false;
      }else if(
          !beforeOperator.includes(expr[i-1]) && 
          !afterOperator.includes(expr[i+1])
        ){
        
        console.log(`Characters surrounding ${c} must be numbers.`)
        return false;
      }
    }

    // Check successive parenthesis
    if(c === "(" && i < n-1 && expr[i+1] === ")"){
      console.log(`Character "${c}" can't be follow by ")"`)
    }

    // Check parenthesis closing order
    if(c === "(") leftParenthesis.push(c);
    if(c === ")") {
      if(leftParenthesis.lenght === 0) return false;
      else leftParenthesis.pop()
    }
  }

  if(leftParenthesis.length > 0) return false;
  
  return true
}

console.log(isCorrect("100*2+6"))
console.log(isCorrect("100*2*12+(1+2)"))
console.log(isCorrect("100*(2+12)/14"))
console.log(isCorrect("3.78+4.78*(1.55)"))