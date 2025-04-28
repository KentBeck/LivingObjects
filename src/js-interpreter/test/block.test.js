const { Interpreter, ast, core } = require('../src');

describe('Block', () => {
  let interpreter;

  beforeEach(() => {
    interpreter = new Interpreter();
  });

  test('should have the correct class', () => {
    const block = new core.STBlock([], [], null);
    expect(block.class).toBe(core.STClass.blockClass);
  });

  test('should capture outer context', () => {
    const outerContext = new core.STContext();
    outerContext.setVariable('x', new core.STInteger(42));
    
    const block = new core.STBlock([], [], outerContext);
    expect(block.outerContext).toBe(outerContext);
  });

  test('should execute with no arguments', () => {
    let executed = false;
    const block = new core.STBlock([], [], null);
    block.value = function() {
      executed = true;
      return new core.STInteger(42);
    };
    
    const result = block.sendMessage('value', [], null);
    expect(executed).toBe(true);
    expect(result.value).toBe(42);
  });

  test('should execute with one argument', () => {
    let capturedArg = null;
    const block = new core.STBlock(['x'], [], null);
    block.value = function(arg) {
      capturedArg = arg;
      return arg;
    };
    
    const arg = new core.STInteger(42);
    const result = block.sendMessage('value:', [arg], null);
    expect(capturedArg).toBe(arg);
    expect(result).toBe(arg);
  });

  test('should execute with two arguments', () => {
    let capturedArg1 = null;
    let capturedArg2 = null;
    const block = new core.STBlock(['x', 'y'], [], null);
    block.value = function(arg1, arg2) {
      capturedArg1 = arg1;
      capturedArg2 = arg2;
      return arg2;
    };
    
    const arg1 = new core.STInteger(42);
    const arg2 = new core.STInteger(24);
    const result = block.sendMessage('value:value:', [arg1, arg2], null);
    expect(capturedArg1).toBe(arg1);
    expect(capturedArg2).toBe(arg2);
    expect(result).toBe(arg2);
  });

  test('should be able to access variables from outer context', () => {
    const outerContext = new core.STContext();
    const x = new core.STInteger(42);
    outerContext.setVariable('x', x);
    
    // Create a block that accesses 'x' from the outer context
    const statements = [
      new ast.Variable('x')
    ];
    
    const block = new core.STBlock([], statements, outerContext);
    
    // Create a context for evaluation
    const evalContext = new core.STContext(outerContext);
    
    // Evaluate the block
    const result = statements[0].evaluate(evalContext);
    expect(result).toBe(x);
  });
});
