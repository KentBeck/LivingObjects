const { Interpreter, ast, core } = require('../src');

describe('Object', () => {
  let interpreter;

  beforeEach(() => {
    interpreter = new Interpreter();
  });

  test('should have a class', () => {
    const object = new core.STObject();
    expect(object.class).toBe(core.STClass.objectClass);
  });

  test('should be able to send messages', () => {
    const object = new core.STObject();
    const result = object.sendMessage('class', [], null);
    expect(result).toBe(core.STClass.objectClass);
  });

  test('should be able to compare equality', () => {
    const object1 = new core.STObject();
    const object2 = new core.STObject();
    
    // Same object should be equal
    expect(object1.equals(object1)).toBe(true);
    
    // Different objects should not be equal
    expect(object1.equals(object2)).toBe(false);
  });

  test('should have a string representation', () => {
    const object = new core.STObject();
    expect(object.toString()).toBe('a Object');
  });

  test('should evaluate equality message send', () => {
    const object1 = new core.STObject();
    const object2 = new core.STObject();
    
    // Test == message
    const result1 = object1.sendMessage('==', [object1], null);
    expect(result1).toBe(core.STBoolean.true);
    
    const result2 = object1.sendMessage('==', [object2], null);
    expect(result2).toBe(core.STBoolean.false);
  });
});
