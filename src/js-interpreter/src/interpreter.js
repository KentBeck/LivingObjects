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
  STUndefinedObject,
  STString,
  STException,
  isSubclassOf,
} = require("./core");

class Interpreter {
  constructor() {
    this.context = new STContext();
    this.initializeGlobals();
    this.initializeBasicMethods();
  }

  // Initialize global variables
  initializeGlobals() {
    this.context.setVariable("Object", STClass.objectClass);
    this.context.setVariable("Class", STClass.classClass);
    this.context.setVariable("Integer", STClass.integerClass);
    this.context.setVariable("Boolean", STClass.booleanClass);
    this.context.setVariable("True", STClass.trueClass);
    this.context.setVariable("False", STClass.falseClass);
    this.context.setVariable("UndefinedObject", STClass.undefinedObjectClass);
    this.context.setVariable("Block", STClass.blockClass);
    this.context.setVariable("String", STClass.stringClass);

    // Exception classes
    this.context.setVariable("Exception", STClass.exceptionClass);
    this.context.setVariable("Error", STClass.errorClass);
    this.context.setVariable(
      "MessageNotUnderstood",
      STClass.messageNotUnderstoodClass
    );
    this.context.setVariable("NotFound", STClass.notFoundClass);
    this.context.setVariable("ZeroDivide", STClass.zeroDivideClass);

    this.context.setVariable("true", STBoolean.true);
    this.context.setVariable("false", STBoolean.false);
    this.context.setVariable("nil", STUndefinedObject.nil);
  }

  // Initialize basic methods for core classes
  initializeBasicMethods() {
    // Object methods
    this.defineMethod(
      STClass.objectClass,
      "==",
      ["anObject"],
      function (self, args) {
        return self.equals(args[0]) ? STBoolean.true : STBoolean.false;
      }
    );

    this.defineMethod(STClass.objectClass, "class", [], function (self) {
      return self.class;
    });

    // Integer methods
    this.defineMethod(
      STClass.integerClass,
      "+",
      ["anInteger"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STInteger) {
          return new STInteger(self.value + other.value);
        }
        throw new Error("Expected an Integer");
      }
    );

    this.defineMethod(
      STClass.integerClass,
      "-",
      ["anInteger"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STInteger) {
          return new STInteger(self.value - other.value);
        }
        throw new Error("Expected an Integer");
      }
    );

    this.defineMethod(
      STClass.integerClass,
      "*",
      ["anInteger"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STInteger) {
          return new STInteger(self.value * other.value);
        }
        throw new Error("Expected an Integer");
      }
    );

    this.defineMethod(
      STClass.integerClass,
      "/",
      ["anInteger"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STInteger) {
          if (other.value === 0) {
            // Create a ZeroDivide exception
            const exception = new STException("Division by zero");
            exception.class = STClass.zeroDivideClass;
            throw exception;
          }
          return new STInteger(Math.floor(self.value / other.value));
        }
        throw new Error("Expected an Integer");
      }
    );

    this.defineMethod(
      STClass.integerClass,
      "<",
      ["anInteger"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STInteger) {
          return self.value < other.value ? STBoolean.true : STBoolean.false;
        }
        throw new Error("Expected an Integer");
      }
    );

    this.defineMethod(
      STClass.integerClass,
      ">",
      ["anInteger"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STInteger) {
          return self.value > other.value ? STBoolean.true : STBoolean.false;
        }
        throw new Error("Expected an Integer");
      }
    );

    this.defineMethod(
      STClass.integerClass,
      "<=",
      ["anInteger"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STInteger) {
          return self.value <= other.value ? STBoolean.true : STBoolean.false;
        }
        throw new Error("Expected an Integer");
      }
    );

    this.defineMethod(
      STClass.integerClass,
      ">=",
      ["anInteger"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STInteger) {
          return self.value >= other.value ? STBoolean.true : STBoolean.false;
        }
        throw new Error("Expected an Integer");
      }
    );

    this.defineMethod(
      STClass.integerClass,
      "=",
      ["anInteger"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STInteger) {
          return self.value === other.value ? STBoolean.true : STBoolean.false;
        }
        return STBoolean.false;
      }
    );

    // Boolean methods
    this.defineMethod(
      STClass.booleanClass,
      "ifTrue:",
      ["aBlock"],
      function (self, args, context) {
        const block = args[0];
        if (self.value) {
          return block.value();
        }
        return STUndefinedObject.nil;
      }
    );

    this.defineMethod(
      STClass.booleanClass,
      "ifFalse:",
      ["aBlock"],
      function (self, args, context) {
        const block = args[0];
        if (!self.value) {
          return block.value();
        }
        return STUndefinedObject.nil;
      }
    );

    this.defineMethod(
      STClass.booleanClass,
      "ifTrue:ifFalse:",
      ["trueBlock", "falseBlock"],
      function (self, args, context) {
        const trueBlock = args[0];
        const falseBlock = args[1];
        return self.value ? trueBlock.value() : falseBlock.value();
      }
    );

    this.defineMethod(STClass.booleanClass, "not", [], function (self) {
      return self.value ? STBoolean.false : STBoolean.true;
    });

    this.defineMethod(
      STClass.booleanClass,
      "&",
      ["aBoolean"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STBoolean) {
          return self.value && other.value ? STBoolean.true : STBoolean.false;
        }
        throw new Error("Expected a Boolean");
      }
    );

    this.defineMethod(
      STClass.booleanClass,
      "|",
      ["aBoolean"],
      function (self, args) {
        const other = args[0];
        if (other instanceof STBoolean) {
          return self.value || other.value ? STBoolean.true : STBoolean.false;
        }
        throw new Error("Expected a Boolean");
      }
    );

    // Block methods
    this.defineMethod(STClass.blockClass, "value", [], function (self) {
      return self.value();
    });

    this.defineMethod(
      STClass.blockClass,
      "value:",
      ["anArg"],
      function (self, args) {
        return self.value(args[0]);
      }
    );

    this.defineMethod(
      STClass.blockClass,
      "value:value:",
      ["firstArg", "secondArg"],
      function (self, args) {
        return self.value(args[0], args[1]);
      }
    );

    // Exception handling methods
    this.defineMethod(
      STClass.blockClass,
      "on:do:",
      ["exceptionClass", "handlerBlock"],
      function (self, args) {
        const exceptionClass = args[0];
        const handlerBlock = args[1];

        try {
          return self.value();
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
    );

    // Exception methods
    this.defineMethod(STClass.exceptionClass, "signal", [], function (self) {
      throw self;
    });

    this.defineMethod(
      STClass.exceptionClass,
      "messageText",
      [],
      function (self) {
        return self.message || "";
      }
    );

    // Class methods for creating new instances
    this.defineMethod(STClass.classClass, "new", [], function (self) {
      return self.newInstance();
    });

    this.defineMethod(
      STClass.classClass,
      "new:super:",
      ["name", "superclass"],
      function (self, args) {
        const name = args[0];
        const superclass = args[1];
        return new STClass(name, superclass);
      }
    );

    this.defineMethod(
      STClass.exceptionClass,
      "new:",
      ["messageText"],
      function (self, args) {
        const messageText = args[0];
        const exception = new STException(messageText);
        exception.class = self;
        return exception;
      }
    );

    // String methods
    this.defineMethod(
      STClass.stringClass,
      ",",
      ["aString"],
      function (self, args) {
        const other = args[0];
        let otherStr = "";

        if (other instanceof STString) {
          otherStr = other.value;
        } else if (other !== undefined && other !== null) {
          otherStr = other.toString();
        }

        return new STString(self.value + otherStr);
      }
    );
  }

  // Helper method to define primitive methods
  defineMethod(classObj, selector, paramNames, implementation) {
    classObj.addMethod(selector, {
      selector,
      parameters: paramNames,
      execute: function (receiver, args, context) {
        return implementation(receiver, args, context);
      },
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

  // Create a new string
  newString(value) {
    return new STString(value);
  }
}

module.exports = Interpreter;
