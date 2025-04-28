/**
 * Exception handling for Smalltalk interpreter
 */

const { STObject, STClass } = require('../core');

// Exception class
class STException extends STObject {
  constructor(message) {
    super();
    this.message = message || '';
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

// Initialize Exception class
function initializeExceptionClasses() {
  // Create the exception class hierarchy
  STClass.exceptionClass = new STClass('Exception', STClass.objectClass);
  STClass.errorClass = new STClass('Error', STClass.exceptionClass);
  STClass.notFoundClass = new STClass('NotFound', STClass.exceptionClass);
  STClass.zeroDivideClass = new STClass('ZeroDivide', STClass.errorClass);
  
  // Set up the class of exception objects
  STClass.exceptionClass.class = STClass.classClass;
  STClass.errorClass.class = STClass.classClass;
  STClass.notFoundClass.class = STClass.classClass;
  STClass.zeroDivideClass.class = STClass.classClass;
}

// Initialize the exception class hierarchy
initializeExceptionClasses();

module.exports = {
  STException
};
