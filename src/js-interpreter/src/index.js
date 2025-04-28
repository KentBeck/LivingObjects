/**
 * Main entry point for the Smalltalk interpreter
 */

const Interpreter = require('./interpreter');
const ast = require('./ast');
const core = require('./core');

module.exports = {
  Interpreter,
  ast,
  core
};
