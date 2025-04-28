/**
 * Core classes for Smalltalk interpreter
 */

const { Block: ASTBlock } = require("./ast");

// Base class for all Smalltalk objects
class STObject {
  constructor() {
    this.class = STClass.objectClass;
  }

  // Send a message to this object
  sendMessage(selector, args, context) {
    try {
      const method = this.class.lookupMethod(selector);
      if (method) {
        return method.execute(this, args, context);
      }

      // If method not found, create a MessageNotUnderstood exception
      const exception = new STException(`Message not understood: ${selector}`);
      exception.class = STClass.messageNotUnderstoodClass;
      exception.receiver = this;
      exception.selector = selector;
      exception.arguments = args;
      throw exception;
    } catch (e) {
      // If it's a Smalltalk exception, propagate it
      if (e instanceof STException) {
        throw e;
      }
      // Otherwise, wrap it in a JavaScript exception
      const exception = new STException(e.message || String(e));
      exception.class = STClass.errorClass;
      exception.jsError = e;
      throw exception;
    }
  }

  // Default implementation of equality
  equals(other) {
    return this === other;
  }

  // Convert to string representation
  toString() {
    return `a ${this.class.name}`;
  }
}

// Class object
class STClass extends STObject {
  constructor(name, superclass, instanceVariables = []) {
    super();
    this.name = name;
    this.superclass = superclass;
    this.instanceVariables = instanceVariables;
    this.methods = new Map();
    this.class = STClass.classClass;
  }

  // Look up a method in this class or its superclasses
  lookupMethod(selector) {
    if (this.methods.has(selector)) {
      return this.methods.get(selector);
    }
    if (this.superclass) {
      return this.superclass.lookupMethod(selector);
    }
    return null;
  }

  // Add a method to this class
  addMethod(selector, method) {
    this.methods.set(selector, method);
  }

  // Create a new instance of this class
  newInstance() {
    const instance = new STObject();
    instance.class = this;
    return instance;
  }
}

// Method object
class STMethod {
  constructor(selector, parameters, body) {
    this.selector = selector;
    this.parameters = parameters;
    this.body = body;
  }

  // Execute this method with the given receiver and arguments
  execute(receiver, args, outerContext) {
    const context = new STContext(outerContext);
    context.setVariable("self", receiver);

    // Bind parameters to arguments
    for (let i = 0; i < this.parameters.length; i++) {
      context.setVariable(this.parameters[i], args[i]);
    }

    // Execute the method body
    return this.body.evaluate(context);
  }
}

// Block object (closure)
class STBlock extends STObject {
  constructor(parameters, statements, outerContext) {
    super();
    this.parameters = parameters;
    this.statements = statements;
    this.outerContext = outerContext;
    this.class = STClass.blockClass;
  }

  // Execute the block with the given arguments
  value(...args) {
    const context = new STContext(this.outerContext);

    // Bind parameters to arguments
    for (let i = 0; i < this.parameters.length; i++) {
      context.setVariable(this.parameters[i], args[i]);
    }

    // Execute the block body
    let result;
    for (const statement of this.statements) {
      result = statement.evaluate(context);
    }
    return result;
  }

  // Execute the block with exception handling
  on_do(exceptionClass, handlerBlock) {
    try {
      return this.value();
    } catch (e) {
      // Check if the exception is an instance of the specified class
      if (
        e instanceof STException &&
        (e.class === exceptionClass ||
          (e.class.superclass && isSubclassOf(e.class, exceptionClass)))
      ) {
        // Call the handler block with the exception as argument
        return handlerBlock.value(e);
      }
      // If not the right type of exception, re-throw it
      throw e;
    }
  }
}

// Execution context
class STContext {
  constructor(parent = null) {
    this.variables = new Map();
    this.parent = parent;
  }

  // Look up a variable in this context or its parent contexts
  lookup(name) {
    if (this.variables.has(name)) {
      return this.variables.get(name);
    }
    if (this.parent) {
      return this.parent.lookup(name);
    }
    throw new Error(`Variable not found: ${name}`);
  }

  // Set a variable in this context
  setVariable(name, value) {
    this.variables.set(name, value);
    return value;
  }

  // Assign a value to a variable in this context or a parent context
  assign(name, value) {
    if (this.variables.has(name)) {
      this.variables.set(name, value);
      return value;
    }
    if (this.parent) {
      return this.parent.assign(name, value);
    }
    // If not found, create it in this context
    return this.setVariable(name, value);
  }

  // Create a new block object
  createBlock(parameters, statements) {
    return new STBlock(parameters, statements, this);
  }
}

// Integer class
class STInteger extends STObject {
  constructor(value) {
    super();
    this.value = value;
    this.class = STClass.integerClass;
  }

  // Override toString
  toString() {
    return this.value.toString();
  }
}

// Boolean class
class STBoolean extends STObject {
  constructor(value) {
    super();
    this.value = value;
    this.class = value ? STClass.trueClass : STClass.falseClass;
  }

  // Override toString
  toString() {
    return this.value.toString();
  }
}

// UndefinedObject class (nil)
class STUndefinedObject extends STObject {
  constructor() {
    super();
    this.class = STClass.undefinedObjectClass;
  }

  // Override toString
  toString() {
    return "nil";
  }
}

// Helper function to check if a class is a subclass of another
function isSubclassOf(classObj, potentialSuperclass) {
  let current = classObj;
  while (current) {
    if (current === potentialSuperclass) {
      return true;
    }
    current = current.superclass;
  }
  return false;
}

// String class
class STString extends STObject {
  constructor(value) {
    super();
    this.value = value || "";
    this.class = STClass.stringClass;
  }

  // Override toString
  toString() {
    return this.value;
  }
}

// Initialize the class hierarchy
function initializeClassHierarchy() {
  // Create the class hierarchy
  STClass.objectClass = new STClass("Object", null);
  STClass.classClass = new STClass("Class", STClass.objectClass);
  STClass.integerClass = new STClass("Integer", STClass.objectClass);
  STClass.booleanClass = new STClass("Boolean", STClass.objectClass);
  STClass.trueClass = new STClass("True", STClass.booleanClass);
  STClass.falseClass = new STClass("False", STClass.booleanClass);
  STClass.undefinedObjectClass = new STClass(
    "UndefinedObject",
    STClass.objectClass
  );
  STClass.blockClass = new STClass("Block", STClass.objectClass);
  STClass.stringClass = new STClass("String", STClass.objectClass);

  // Exception classes
  STClass.exceptionClass = new STClass("Exception", STClass.objectClass);
  STClass.errorClass = new STClass("Error", STClass.exceptionClass);
  STClass.messageNotUnderstoodClass = new STClass(
    "MessageNotUnderstood",
    STClass.exceptionClass
  );
  STClass.notFoundClass = new STClass("NotFound", STClass.exceptionClass);
  STClass.zeroDivideClass = new STClass("ZeroDivide", STClass.errorClass);

  // Set up the class of class objects
  STClass.objectClass.class = STClass.classClass;
  STClass.classClass.class = STClass.classClass;
  STClass.integerClass.class = STClass.classClass;
  STClass.booleanClass.class = STClass.classClass;
  STClass.trueClass.class = STClass.classClass;
  STClass.falseClass.class = STClass.classClass;
  STClass.undefinedObjectClass.class = STClass.classClass;
  STClass.blockClass.class = STClass.classClass;
  STClass.stringClass.class = STClass.classClass;

  // Set up the class of exception objects
  STClass.exceptionClass.class = STClass.classClass;
  STClass.errorClass.class = STClass.classClass;
  STClass.messageNotUnderstoodClass.class = STClass.classClass;
  STClass.notFoundClass.class = STClass.classClass;
  STClass.zeroDivideClass.class = STClass.classClass;

  // Create singleton instances
  STBoolean.true = new STBoolean(true);
  STBoolean.false = new STBoolean(false);
  STUndefinedObject.nil = new STUndefinedObject();
}

// Initialize the class hierarchy
initializeClassHierarchy();

// Exception class
class STException extends STObject {
  constructor(message) {
    super();
    this.message = message || "";
    this.class = STClass.exceptionClass;
  }

  // Override toString
  toString() {
    return `${this.class.name}: ${this.message}`;
  }

  // Signal this exception
  signal() {
    throw this;
  }
}

module.exports = {
  STObject,
  STClass,
  STMethod,
  STBlock,
  STContext,
  STInteger,
  STBoolean,
  STUndefinedObject,
  STString,
  STException,
  isSubclassOf,
};
