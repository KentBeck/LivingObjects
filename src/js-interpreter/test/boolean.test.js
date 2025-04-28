const { Interpreter, ast, core } = require('../src');

describe('Boolean', () => {
  let interpreter;

  beforeEach(() => {
    interpreter = new Interpreter();
  });

  test('should have true and false singletons', () => {
    expect(core.STBoolean.true.value).toBe(true);
    expect(core.STBoolean.false.value).toBe(false);
  });

  test('should have appropriate classes', () => {
    expect(core.STBoolean.true.class).toBe(core.STClass.trueClass);
    expect(core.STBoolean.false.class).toBe(core.STClass.falseClass);
  });

  test('should have a string representation', () => {
    expect(core.STBoolean.true.toString()).toBe('true');
    expect(core.STBoolean.false.toString()).toBe('false');
  });

  test('should perform logical not', () => {
    const result1 = core.STBoolean.true.sendMessage('not', [], null);
    expect(result1).toBe(core.STBoolean.false);
    
    const result2 = core.STBoolean.false.sendMessage('not', [], null);
    expect(result2).toBe(core.STBoolean.true);
  });

  test('should perform logical and', () => {
    const result1 = core.STBoolean.true.sendMessage('&', [core.STBoolean.true], null);
    expect(result1).toBe(core.STBoolean.true);
    
    const result2 = core.STBoolean.true.sendMessage('&', [core.STBoolean.false], null);
    expect(result2).toBe(core.STBoolean.false);
    
    const result3 = core.STBoolean.false.sendMessage('&', [core.STBoolean.true], null);
    expect(result3).toBe(core.STBoolean.false);
    
    const result4 = core.STBoolean.false.sendMessage('&', [core.STBoolean.false], null);
    expect(result4).toBe(core.STBoolean.false);
  });

  test('should perform logical or', () => {
    const result1 = core.STBoolean.true.sendMessage('|', [core.STBoolean.true], null);
    expect(result1).toBe(core.STBoolean.true);
    
    const result2 = core.STBoolean.true.sendMessage('|', [core.STBoolean.false], null);
    expect(result2).toBe(core.STBoolean.true);
    
    const result3 = core.STBoolean.false.sendMessage('|', [core.STBoolean.true], null);
    expect(result3).toBe(core.STBoolean.true);
    
    const result4 = core.STBoolean.false.sendMessage('|', [core.STBoolean.false], null);
    expect(result4).toBe(core.STBoolean.false);
  });

  test('should execute ifTrue: block', () => {
    let executed = false;
    const block = new core.STBlock([], [], null);
    block.value = function() {
      executed = true;
      return new core.STInteger(42);
    };
    
    const result1 = core.STBoolean.true.sendMessage('ifTrue:', [block], null);
    expect(executed).toBe(true);
    expect(result1.value).toBe(42);
    
    executed = false;
    const result2 = core.STBoolean.false.sendMessage('ifTrue:', [block], null);
    expect(executed).toBe(false);
    expect(result2).toBe(core.STUndefinedObject.nil);
  });

  test('should execute ifFalse: block', () => {
    let executed = false;
    const block = new core.STBlock([], [], null);
    block.value = function() {
      executed = true;
      return new core.STInteger(42);
    };
    
    const result1 = core.STBoolean.false.sendMessage('ifFalse:', [block], null);
    expect(executed).toBe(true);
    expect(result1.value).toBe(42);
    
    executed = false;
    const result2 = core.STBoolean.true.sendMessage('ifFalse:', [block], null);
    expect(executed).toBe(false);
    expect(result2).toBe(core.STUndefinedObject.nil);
  });

  test('should execute ifTrue:ifFalse: blocks', () => {
    let executedTrue = false;
    let executedFalse = false;
    
    const trueBlock = new core.STBlock([], [], null);
    trueBlock.value = function() {
      executedTrue = true;
      return new core.STInteger(42);
    };
    
    const falseBlock = new core.STBlock([], [], null);
    falseBlock.value = function() {
      executedFalse = true;
      return new core.STInteger(24);
    };
    
    const result1 = core.STBoolean.true.sendMessage('ifTrue:ifFalse:', [trueBlock, falseBlock], null);
    expect(executedTrue).toBe(true);
    expect(executedFalse).toBe(false);
    expect(result1.value).toBe(42);
    
    executedTrue = false;
    executedFalse = false;
    
    const result2 = core.STBoolean.false.sendMessage('ifTrue:ifFalse:', [trueBlock, falseBlock], null);
    expect(executedTrue).toBe(false);
    expect(executedFalse).toBe(true);
    expect(result2.value).toBe(24);
  });
});
