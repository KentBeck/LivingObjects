/**
 * Example demonstrating exception handling in the Smalltalk interpreter
 */

const { Interpreter, ast, core } = require("../src");

// Create an interpreter instance
const interpreter = new Interpreter();

// Create a program that demonstrates exception handling
function createExceptionHandlingDemo() {
  // This program demonstrates a factorial function with error handling
  // factorial := [:n |
  //   (n < 0)
  //     ifTrue: [
  //       (Error new: 'Factorial not defined for negative numbers') signal
  //     ]
  //     ifFalse: [
  //       n = 0
  //         ifTrue: [1]
  //         ifFalse: [n * (factorial value: n - 1)]
  //     ]
  // ].
  //
  // result := [
  //   factorial value: -5
  // ] on: Error do: [:ex | 'Error: ', ex messageText].

  // First, create the factorial block
  const nVar = new ast.Variable("n");
  const zero = new ast.Literal(new core.STInteger(0));
  const one = new ast.Literal(new core.STInteger(1));
  const negativeOne = new ast.Literal(new core.STInteger(-1));

  // n < 0
  const isNegative = new ast.MessageSend(nVar, "<", [zero]);

  // Create an error for negative numbers
  const errorMessage = "Factorial not defined for negative numbers";
  const exception = new core.STException(errorMessage);
  exception.class = core.STClass.errorClass;
  const createError = new ast.MessageSend(
    new ast.Literal(exception),
    "signal",
    []
  );

  // [Error new: '...'] signal
  const signalErrorBlock = new ast.Block([], [createError]);

  // n = 0
  const isZero = new ast.MessageSend(nVar, "=", [zero]);

  // n - 1
  const nMinusOne = new ast.MessageSend(nVar, "-", [one]);

  // factorial variable
  const factorialVar = new ast.Variable("factorial");

  // factorial value: n - 1
  const recursiveCall = new ast.MessageSend(factorialVar, "value:", [
    nMinusOne,
  ]);

  // n * (factorial value: n - 1)
  const multiply = new ast.MessageSend(nVar, "*", [recursiveCall]);

  // [n * (factorial value: n - 1)]
  const recursiveBlock = new ast.Block([], [multiply]);

  // [1]
  const returnOneBlock = new ast.Block([], [one]);

  // n = 0 ifTrue: [1] ifFalse: [n * (factorial value: n - 1)]
  const zeroCase = new ast.MessageSend(isZero, "ifTrue:ifFalse:", [
    returnOneBlock,
    recursiveBlock,
  ]);

  // [n = 0 ifTrue: [1] ifFalse: [n * (factorial value: n - 1)]]
  const notNegativeBlock = new ast.Block([], [zeroCase]);

  // (n < 0) ifTrue: [signal error] ifFalse: [...]
  const checkNegative = new ast.MessageSend(isNegative, "ifTrue:ifFalse:", [
    signalErrorBlock,
    notNegativeBlock,
  ]);

  // [:n | ... ]
  const factorialBlock = new ast.Block(["n"], [checkNegative]);

  // factorial := [:n | ... ]
  const assignFactorial = new ast.Assignment(factorialVar, factorialBlock);

  // Create a call to factorial with a negative number
  const minusFive = new ast.Literal(new core.STInteger(-5));
  const callWithNegative = new ast.MessageSend(factorialVar, "value:", [
    minusFive,
  ]);

  // [factorial value: -5]
  const tryBlock = new ast.Block([], [callWithNegative]);

  // Create a handler for Error
  const exVar = new ast.Variable("ex");
  const getMessageText = new ast.MessageSend(exVar, "messageText", []);
  const errorPrefix = new ast.Literal(interpreter.newString("Error: "));
  const concatenate = new ast.MessageSend(errorPrefix, ",", [getMessageText]);
  const handlerBlock = new ast.Block(["ex"], [concatenate]);

  // [factorial value: -5] on: Error do: [:ex | 'Error: ', ex messageText]
  const onDo = new ast.MessageSend(tryBlock, "on:do:", [
    new ast.Literal(core.STClass.errorClass),
    handlerBlock,
  ]);

  // Create a program
  return new ast.Program([
    // Define the factorial function and call it with error handling
    assignFactorial,
    onDo,
  ]);
}

// Create and evaluate the program
const program = createExceptionHandlingDemo();
const result = interpreter.evaluate(program);

console.log("Result:", result);

// Now try with a valid input
function createValidFactorialDemo() {
  // factorial := [:n |
  //   (n < 0)
  //     ifTrue: [
  //       (Error new: 'Factorial not defined for negative numbers') signal
  //     ]
  //     ifFalse: [
  //       n = 0
  //         ifTrue: [1]
  //         ifFalse: [n * (factorial value: n - 1)]
  //     ]
  // ].
  //
  // result := factorial value: 5.

  // Reuse the factorial function from above
  const factorialVar = new ast.Variable("factorial");
  const five = new ast.Literal(new core.STInteger(5));
  const callWithFive = new ast.MessageSend(factorialVar, "value:", [five]);

  // Create a program
  return new ast.Program([callWithFive]);
}

// Create and evaluate the program with a valid input
const validProgram = createValidFactorialDemo();
const validResult = interpreter.evaluate(validProgram);

console.log("Valid Result:", validResult);
