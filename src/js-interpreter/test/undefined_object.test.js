const { Interpreter, ast, core } = require('../src');

describe('UndefinedObject', () => {
  let interpreter;

  beforeEach(() => {
    interpreter = new Interpreter();
  });

  test('should have a nil singleton', () => {
    expect(core.STUndefinedObject.nil).toBeInstanceOf(core.STUndefinedObject);
  });

  test('should have the correct class', () => {
    expect(core.STUndefinedObject.nil.class).toBe(core.STClass.undefinedObjectClass);
  });

  test('should have a string representation', () => {
    expect(core.STUndefinedObject.nil.toString()).toBe('nil');
  });

  test('should be equal to itself', () => {
    const result = core.STUndefinedObject.nil.sendMessage('==', [core.STUndefinedObject.nil], null);
    expect(result).toBe(core.STBoolean.true);
  });

  test('should not be equal to other objects', () => {
    const object = new core.STObject();
    const result = core.STUndefinedObject.nil.sendMessage('==', [object], null);
    expect(result).toBe(core.STBoolean.false);
  });
});
