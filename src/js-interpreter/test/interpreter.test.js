const { Interpreter, ast, core } = require('../src');

describe('Interpreter', () => {
  let interpreter;

  beforeEach(() => {
    interpreter = new Interpreter();
  });

  test('should evaluate literals', () => {
    const literal = new ast.Literal(new core.STInteger(42));
    const result = interpreter.evaluate(literal);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(42);
  });

  test('should evaluate variables', () => {
    // Set a variable in the interpreter context
    interpreter.context.setVariable('x', new core.STInteger(42));
    
    // Create a variable node
    const variable = new ast.Variable('x');
    
    // Evaluate the variable
    const result = interpreter.evaluate(variable);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(42);
  });

  test('should evaluate assignments', () => {
    // Create an assignment node
    const variable = new ast.Variable('x');
    const value = new ast.Literal(new core.STInteger(42));
    const assignment = new ast.Assignment(variable, value);
    
    // Evaluate the assignment
    const result = interpreter.evaluate(assignment);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(42);
    
    // Check that the variable was set
    const x = interpreter.context.lookup('x');
    expect(x).toBeInstanceOf(core.STInteger);
    expect(x.value).toBe(42);
  });

  test('should evaluate message sends', () => {
    // Create a message send: 3 + 4
    const receiver = new ast.Literal(new core.STInteger(3));
    const arg = new ast.Literal(new core.STInteger(4));
    const messageSend = new ast.MessageSend(receiver, '+', [arg]);
    
    // Evaluate the message send
    const result = interpreter.evaluate(messageSend);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(7);
  });

  test('should evaluate blocks', () => {
    // Create a block: [x + 1]
    const xVar = new ast.Variable('x');
    const one = new ast.Literal(new core.STInteger(1));
    const addition = new ast.MessageSend(xVar, '+', [one]);
    const block = new ast.Block(['x'], [addition]);
    
    // Evaluate the block to get a block object
    const blockObj = interpreter.evaluate(block);
    expect(blockObj).toBeInstanceOf(core.STBlock);
    
    // Execute the block with an argument
    const arg = new core.STInteger(5);
    const result = blockObj.value(arg);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(6);
  });

  test('should evaluate conditionals', () => {
    // Create a conditional: true ifTrue: [42] ifFalse: [24]
    const trueObj = new ast.Literal(core.STBoolean.true);
    
    // Create the true block: [42]
    const trueValue = new ast.Literal(new core.STInteger(42));
    const trueBlock = new ast.Block([], [trueValue]);
    
    // Create the false block: [24]
    const falseValue = new ast.Literal(new core.STInteger(24));
    const falseBlock = new ast.Block([], [falseValue]);
    
    // Create the message send: true ifTrue:ifFalse: [42] [24]
    const conditional = new ast.MessageSend(trueObj, 'ifTrue:ifFalse:', [trueBlock, falseBlock]);
    
    // Evaluate the conditional
    const result = interpreter.evaluate(conditional);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(42);
  });

  test('should evaluate a program', () => {
    // Create a program:
    // x := 3.
    // y := 4.
    // x + y
    
    // x := 3
    const xVar = new ast.Variable('x');
    const three = new ast.Literal(new core.STInteger(3));
    const assignX = new ast.Assignment(xVar, three);
    
    // y := 4
    const yVar = new ast.Variable('y');
    const four = new ast.Literal(new core.STInteger(4));
    const assignY = new ast.Assignment(yVar, four);
    
    // x + y
    const addition = new ast.MessageSend(xVar, '+', [yVar]);
    
    // Create the program
    const program = new ast.Program([assignX, assignY, addition]);
    
    // Evaluate the program
    const result = interpreter.evaluate(program);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(7);
  });

  test('should evaluate a factorial function', () => {
    // Create a factorial function:
    // factorial := [:n | n = 0 ifTrue: [1] ifFalse: [n * (factorial value: n - 1)]].
    // factorial value: 5
    
    // First, create the factorial block
    const nVar = new ast.Variable('n');
    const zero = new ast.Literal(new core.STInteger(0));
    const one = new ast.Literal(new core.STInteger(1));
    
    // n = 0
    const equals = new ast.MessageSend(nVar, '=', [zero]);
    
    // [1] (true block)
    const trueBlock = new ast.Block([], [one]);
    
    // n - 1
    const minusOne = new ast.MessageSend(nVar, '-', [one]);
    
    // factorial variable
    const factorialVar = new ast.Variable('factorial');
    
    // factorial value: n - 1
    const recursiveCall = new ast.MessageSend(factorialVar, 'value:', [minusOne]);
    
    // n * (factorial value: n - 1)
    const multiply = new ast.MessageSend(nVar, '*', [recursiveCall]);
    
    // [n * (factorial value: n - 1)] (false block)
    const falseBlock = new ast.Block([], [multiply]);
    
    // n = 0 ifTrue: [1] ifFalse: [n * (factorial value: n - 1)]
    const conditional = new ast.MessageSend(equals, 'ifTrue:ifFalse:', [trueBlock, falseBlock]);
    
    // [:n | n = 0 ifTrue: [1] ifFalse: [n * (factorial value: n - 1)]]
    const factorialBlock = new ast.Block(['n'], [conditional]);
    
    // factorial := [:n | n = 0 ifTrue: [1] ifFalse: [n * (factorial value: n - 1)]]
    const assignFactorial = new ast.Assignment(factorialVar, factorialBlock);
    
    // Create the argument for the factorial function
    const five = new ast.Literal(new core.STInteger(5));
    
    // factorial value: 5
    const callFactorial = new ast.MessageSend(factorialVar, 'value:', [five]);
    
    // Create the program
    const program = new ast.Program([assignFactorial, callFactorial]);
    
    // Evaluate the program
    const result = interpreter.evaluate(program);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(120); // 5! = 120
  });
});
