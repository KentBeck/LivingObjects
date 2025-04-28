/**
 * Smalltalk AST-based interpreter
 */

const {
  STObject,
  STClass,
  STMethod,
  STBlock,
  STContext,
  STInteger,
  STBoolean,
  STUndefinedObject
} = require('./core');

class Interpreter {
  constructor() {
    this.context = new STContext();
    this.initializeGlobals();
    this.initializeBasicMethods();
  }

  // Initialize global variables
  initializeGlobals() {
    this.context.setVariable('Object', STClass.objectClass);
    this.context.setVariable('Class', STClass.classClass);
    this.context.setVariable('Integer', STClass.integerClass);
    this.context.setVariable('Boolean', STClass.booleanClass);
    this.context.setVariable('True', STClass.trueClass);
    this.context.setVariable('False', STClass.falseClass);
    this.context.setVariable('UndefinedObject', STClass.undefinedObjectClass);
    this.context.setVariable('Block', STClass.blockClass);
    
    this.context.setVariable('true', STBoolean.true);
    this.context.setVariable('false', STBoolean.false);
    this.context.setVariable('nil', STUndefinedObject.nil);
  }

  // Initialize basic methods for core classes
  initializeBasicMethods() {
    // Object methods
    this.defineMethod(STClass.objectClass, '==', ['anObject'], function(self, args) {
      return self.equals(args[0]) ? STBoolean.true : STBoolean.false;
    });

    this.defineMethod(STClass.objectClass, 'class', [], function(self) {
      return self.class;
    });

    // Integer methods
    this.defineMethod(STClass.integerClass, '+', ['anInteger'], function(self, args) {
      const other = args[0];
      if (other instanceof STInteger) {
        return new STInteger(self.value + other.value);
      }
      throw new Error('Expected an Integer');
    });

    this.defineMethod(STClass.integerClass, '-', ['anInteger'], function(self, args) {
      const other = args[0];
      if (other instanceof STInteger) {
        return new STInteger(self.value - other.value);
      }
      throw new Error('Expected an Integer');
    });

    this.defineMethod(STClass.integerClass, '*', ['anInteger'], function(self, args) {
      const other = args[0];
      if (other instanceof STInteger) {
        return new STInteger(self.value * other.value);
      }
      throw new Error('Expected an Integer');
    });

    this.defineMethod(STClass.integerClass, '/', ['anInteger'], function(self, args) {
      const other = args[0];
      if (other instanceof STInteger) {
        if (other.value === 0) {
          throw new Error('Division by zero');
        }
        return new STInteger(Math.floor(self.value / other.value));
      }
      throw new Error('Expected an Integer');
    });

    this.defineMethod(STClass.integerClass, '<', ['anInteger'], function(self, args) {
      const other = args[0];
      if (other instanceof STInteger) {
        return self.value < other.value ? STBoolean.true : STBoolean.false;
      }
      throw new Error('Expected an Integer');
    });

    this.defineMethod(STClass.integerClass, '>', ['anInteger'], function(self, args) {
      const other = args[0];
      if (other instanceof STInteger) {
        return self.value > other.value ? STBoolean.true : STBoolean.false;
      }
      throw new Error('Expected an Integer');
    });

    this.defineMethod(STClass.integerClass, '<=', ['anInteger'], function(self, args) {
      const other = args[0];
      if (other instanceof STInteger) {
        return self.value <= other.value ? STBoolean.true : STBoolean.false;
      }
      throw new Error('Expected an Integer');
    });

    this.defineMethod(STClass.integerClass, '>=', ['anInteger'], function(self, args) {
      const other = args[0];
      if (other instanceof STInteger) {
        return self.value >= other.value ? STBoolean.true : STBoolean.false;
      }
      throw new Error('Expected an Integer');
    });

    this.defineMethod(STClass.integerClass, '=', ['anInteger'], function(self, args) {
      const other = args[0];
      if (other instanceof STInteger) {
        return self.value === other.value ? STBoolean.true : STBoolean.false;
      }
      return STBoolean.false;
    });

    // Boolean methods
    this.defineMethod(STClass.booleanClass, 'ifTrue:', ['aBlock'], function(self, args, context) {
      const block = args[0];
      if (self.value) {
        return block.value();
      }
      return STUndefinedObject.nil;
    });

    this.defineMethod(STClass.booleanClass, 'ifFalse:', ['aBlock'], function(self, args, context) {
      const block = args[0];
      if (!self.value) {
        return block.value();
      }
      return STUndefinedObject.nil;
    });

    this.defineMethod(STClass.booleanClass, 'ifTrue:ifFalse:', ['trueBlock', 'falseBlock'], function(self, args, context) {
      const trueBlock = args[0];
      const falseBlock = args[1];
      return self.value ? trueBlock.value() : falseBlock.value();
    });

    this.defineMethod(STClass.booleanClass, 'not', [], function(self) {
      return self.value ? STBoolean.false : STBoolean.true;
    });

    this.defineMethod(STClass.booleanClass, '&', ['aBoolean'], function(self, args) {
      const other = args[0];
      if (other instanceof STBoolean) {
        return self.value && other.value ? STBoolean.true : STBoolean.false;
      }
      throw new Error('Expected a Boolean');
    });

    this.defineMethod(STClass.booleanClass, '|', ['aBoolean'], function(self, args) {
      const other = args[0];
      if (other instanceof STBoolean) {
        return self.value || other.value ? STBoolean.true : STBoolean.false;
      }
      throw new Error('Expected a Boolean');
    });

    // Block methods
    this.defineMethod(STClass.blockClass, 'value', [], function(self) {
      return self.value();
    });

    this.defineMethod(STClass.blockClass, 'value:', ['anArg'], function(self, args) {
      return self.value(args[0]);
    });

    this.defineMethod(STClass.blockClass, 'value:value:', ['firstArg', 'secondArg'], function(self, args) {
      return self.value(args[0], args[1]);
    });
  }

  // Helper method to define primitive methods
  defineMethod(classObj, selector, paramNames, implementation) {
    classObj.addMethod(selector, {
      selector,
      parameters: paramNames,
      execute: function(receiver, args, context) {
        return implementation(receiver, args, context);
      }
    });
  }

  // Evaluate an AST
  evaluate(ast) {
    return ast.evaluate(this.context);
  }

  // Create a new integer
  newInteger(value) {
    return new STInteger(value);
  }

  // Create a new boolean
  newBoolean(value) {
    return value ? STBoolean.true : STBoolean.false;
  }

  // Get the nil object
  nil() {
    return STUndefinedObject.nil;
  }
}

module.exports = Interpreter;
