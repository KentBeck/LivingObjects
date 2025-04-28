const { Interpreter, ast, core } = require('../src');

describe('Exception', () => {
  let interpreter;

  beforeEach(() => {
    interpreter = new Interpreter();
  });

  test('should have the correct class hierarchy', () => {
    expect(core.STClass.exceptionClass.superclass).toBe(core.STClass.objectClass);
    expect(core.STClass.errorClass.superclass).toBe(core.STClass.exceptionClass);
    expect(core.STClass.zeroDivideClass.superclass).toBe(core.STClass.errorClass);
    expect(core.STClass.messageNotUnderstoodClass.superclass).toBe(core.STClass.exceptionClass);
  });

  test('should be able to create and signal exceptions', () => {
    const exception = new core.STException('Test exception');
    expect(exception.class).toBe(core.STClass.exceptionClass);
    expect(exception.message).toBe('Test exception');
    
    // Test signaling an exception
    expect(() => {
      exception.signal();
    }).toThrow(core.STException);
  });

  test('should handle division by zero with ZeroDivide exception', () => {
    // Create a division by zero operation: 5 / 0
    const five = new ast.Literal(new core.STInteger(5));
    const zero = new ast.Literal(new core.STInteger(0));
    const division = new ast.MessageSend(five, '/', [zero]);
    
    // Evaluate the division without a handler - should throw
    expect(() => {
      interpreter.evaluate(division);
    }).toThrow(core.STException);
    
    // Now create a handler for ZeroDivide
    const handler = new ast.Block(['ex'], [
      new ast.Literal(new core.STInteger(999)) // Return 999 on division by zero
    ]);
    
    // Create a block that performs the division
    const divisionBlock = new ast.Block([], [division]);
    
    // Create a message send: [5 / 0] on: ZeroDivide do: [:ex | 999]
    const onDo = new ast.MessageSend(
      divisionBlock,
      'on:do:',
      [new ast.Literal(core.STClass.zeroDivideClass), handler]
    );
    
    // Evaluate with the handler
    const result = interpreter.evaluate(onDo);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(999);
  });

  test('should handle message not understood', () => {
    // Create a message send to an object with a non-existent method
    const obj = new core.STObject();
    const messageSend = new ast.MessageSend(
      new ast.Literal(obj),
      'nonExistentMethod',
      []
    );
    
    // Create a handler for MessageNotUnderstood
    const handler = new ast.Block(['ex'], [
      new ast.Literal(new core.STInteger(42)) // Return 42 on message not understood
    ]);
    
    // Create a block that performs the message send
    const messageBlock = new ast.Block([], [messageSend]);
    
    // Create a message send: [obj nonExistentMethod] on: MessageNotUnderstood do: [:ex | 42]
    const onDo = new ast.MessageSend(
      messageBlock,
      'on:do:',
      [new ast.Literal(core.STClass.messageNotUnderstoodClass), handler]
    );
    
    // Evaluate with the handler
    const result = interpreter.evaluate(onDo);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(42);
  });

  test('should only catch exceptions of the specified class', () => {
    // Create a division by zero operation: 5 / 0
    const five = new ast.Literal(new core.STInteger(5));
    const zero = new ast.Literal(new core.STInteger(0));
    const division = new ast.MessageSend(five, '/', [zero]);
    
    // Create a handler for NotFound (wrong exception type)
    const wrongHandler = new ast.Block(['ex'], [
      new ast.Literal(new core.STInteger(999))
    ]);
    
    // Create a block that performs the division
    const divisionBlock = new ast.Block([], [division]);
    
    // Create a message send: [5 / 0] on: NotFound do: [:ex | 999]
    const onDoWrong = new ast.MessageSend(
      divisionBlock,
      'on:do:',
      [new ast.Literal(core.STClass.notFoundClass), wrongHandler]
    );
    
    // Evaluate with the wrong handler - should still throw
    expect(() => {
      interpreter.evaluate(onDoWrong);
    }).toThrow(core.STException);
    
    // Now create a handler for Exception (parent class of ZeroDivide)
    const correctHandler = new ast.Block(['ex'], [
      new ast.Literal(new core.STInteger(888))
    ]);
    
    // Create a message send: [5 / 0] on: Exception do: [:ex | 888]
    const onDoCorrect = new ast.MessageSend(
      divisionBlock,
      'on:do:',
      [new ast.Literal(core.STClass.exceptionClass), correctHandler]
    );
    
    // Evaluate with the correct handler
    const result = interpreter.evaluate(onDoCorrect);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(888);
  });

  test('should be able to access exception properties in handler', () => {
    // Create a custom exception
    const exception = new core.STException('Custom error message');
    exception.class = core.STClass.errorClass;
    
    // Create a block that signals the exception
    const signalBlock = new ast.Block([], [
      new ast.MessageSend(new ast.Literal(exception), 'signal', [])
    ]);
    
    // Create a handler that returns the exception's message
    const handler = new ast.Block(['ex'], [
      new ast.MessageSend(new ast.Variable('ex'), 'messageText', [])
    ]);
    
    // Create a message send: [exception signal] on: Error do: [:ex | ex messageText]
    const onDo = new ast.MessageSend(
      signalBlock,
      'on:do:',
      [new ast.Literal(core.STClass.errorClass), handler]
    );
    
    // Evaluate with the handler
    const result = interpreter.evaluate(onDo);
    expect(result).toBe('Custom error message');
  });
});
