const { Interpreter, ast, core } = require('../src');

describe('Class', () => {
  let interpreter;

  beforeEach(() => {
    interpreter = new Interpreter();
  });

  test('should have a name', () => {
    const objectClass = core.STClass.objectClass;
    expect(objectClass.name).toBe('Object');
  });

  test('should have a class', () => {
    const objectClass = core.STClass.objectClass;
    expect(objectClass.class).toBe(core.STClass.classClass);
  });

  test('should be able to create instances', () => {
    const objectClass = core.STClass.objectClass;
    const instance = objectClass.newInstance();
    expect(instance.class).toBe(objectClass);
  });

  test('should be able to add and lookup methods', () => {
    const testClass = new core.STClass('TestClass', core.STClass.objectClass);
    
    // Add a method
    const method = {
      selector: 'testMethod',
      parameters: [],
      execute: function() {
        return 'test result';
      }
    };
    
    testClass.addMethod('testMethod', method);
    
    // Look up the method
    const foundMethod = testClass.lookupMethod('testMethod');
    expect(foundMethod).toBe(method);
  });

  test('should inherit methods from superclass', () => {
    const superClass = new core.STClass('SuperClass', core.STClass.objectClass);
    const subClass = new core.STClass('SubClass', superClass);
    
    // Add a method to the superclass
    const method = {
      selector: 'inheritedMethod',
      parameters: [],
      execute: function() {
        return 'inherited result';
      }
    };
    
    superClass.addMethod('inheritedMethod', method);
    
    // Look up the method from the subclass
    const foundMethod = subClass.lookupMethod('inheritedMethod');
    expect(foundMethod).toBe(method);
  });
});
