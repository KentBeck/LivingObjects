/**
 * AST classes for Smalltalk interpreter
 */

class Node {
  constructor() {}
}

class Program extends Node {
  constructor(statements) {
    super();
    this.statements = statements || [];
  }

  evaluate(context) {
    let result;
    for (const statement of this.statements) {
      result = statement.evaluate(context);
    }
    return result;
  }
}

class Literal extends Node {
  constructor(value) {
    super();
    this.value = value;
  }

  evaluate(context) {
    return this.value;
  }
}

class Variable extends Node {
  constructor(name) {
    super();
    this.name = name;
  }

  evaluate(context) {
    return context.lookup(this.name);
  }
}

class Assignment extends Node {
  constructor(variable, expression) {
    super();
    this.variable = variable;
    this.expression = expression;
  }

  evaluate(context) {
    const value = this.expression.evaluate(context);
    context.assign(this.variable.name, value);
    return value;
  }
}

class MessageSend extends Node {
  constructor(receiver, selector, args) {
    super();
    this.receiver = receiver;
    this.selector = selector;
    this.args = args || [];
  }

  evaluate(context) {
    const receiver = this.receiver.evaluate(context);
    const args = this.args.map(arg => arg.evaluate(context));
    return receiver.sendMessage(this.selector, args, context);
  }
}

class Block extends Node {
  constructor(parameters, statements) {
    super();
    this.parameters = parameters || [];
    this.statements = statements || [];
  }

  evaluate(context) {
    // Create a Block object that captures the current context
    return context.createBlock(this.parameters, this.statements, context);
  }
}

module.exports = {
  Node,
  Program,
  Literal,
  Variable,
  Assignment,
  MessageSend,
  Block
};
