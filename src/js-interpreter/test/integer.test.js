const { Interpreter, ast, core } = require('../src');

describe('Integer', () => {
  let interpreter;

  beforeEach(() => {
    interpreter = new Interpreter();
  });

  test('should have a value', () => {
    const integer = new core.STInteger(42);
    expect(integer.value).toBe(42);
  });

  test('should have a class', () => {
    const integer = new core.STInteger(42);
    expect(integer.class).toBe(core.STClass.integerClass);
  });

  test('should have a string representation', () => {
    const integer = new core.STInteger(42);
    expect(integer.toString()).toBe('42');
  });

  test('should perform addition', () => {
    const int1 = new core.STInteger(5);
    const int2 = new core.STInteger(3);
    
    const result = int1.sendMessage('+', [int2], null);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(8);
  });

  test('should perform subtraction', () => {
    const int1 = new core.STInteger(5);
    const int2 = new core.STInteger(3);
    
    const result = int1.sendMessage('-', [int2], null);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(2);
  });

  test('should perform multiplication', () => {
    const int1 = new core.STInteger(5);
    const int2 = new core.STInteger(3);
    
    const result = int1.sendMessage('*', [int2], null);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(15);
  });

  test('should perform division', () => {
    const int1 = new core.STInteger(10);
    const int2 = new core.STInteger(3);
    
    const result = int1.sendMessage('/', [int2], null);
    expect(result).toBeInstanceOf(core.STInteger);
    expect(result.value).toBe(3); // Integer division
  });

  test('should compare less than', () => {
    const int1 = new core.STInteger(3);
    const int2 = new core.STInteger(5);
    
    const result = int1.sendMessage('<', [int2], null);
    expect(result).toBe(core.STBoolean.true);
    
    const result2 = int2.sendMessage('<', [int1], null);
    expect(result2).toBe(core.STBoolean.false);
  });

  test('should compare greater than', () => {
    const int1 = new core.STInteger(5);
    const int2 = new core.STInteger(3);
    
    const result = int1.sendMessage('>', [int2], null);
    expect(result).toBe(core.STBoolean.true);
    
    const result2 = int2.sendMessage('>', [int1], null);
    expect(result2).toBe(core.STBoolean.false);
  });

  test('should compare equality', () => {
    const int1 = new core.STInteger(5);
    const int2 = new core.STInteger(5);
    const int3 = new core.STInteger(3);
    
    const result1 = int1.sendMessage('=', [int2], null);
    expect(result1).toBe(core.STBoolean.true);
    
    const result2 = int1.sendMessage('=', [int3], null);
    expect(result2).toBe(core.STBoolean.false);
  });
});
