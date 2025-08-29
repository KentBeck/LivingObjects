#pragma once

#include "object.h"
#include "tagged_value.h"
#include <memory>
#include <string>
#include <vector>

namespace smalltalk {

// Forward declarations
class Interpreter;
struct MethodContext;

/**
 * Base class for all Smalltalk exceptions
 * Represents the Exception class in Smalltalk
 */
class SmalltalkException : public Object {
public:
  SmalltalkException(const std::string &message = "",
                     const std::string &exceptionClass = "Exception");
  virtual ~SmalltalkException() = default;

  // Exception interface
  virtual std::string getMessage() const { return message_; }
  virtual std::string getExceptionClass() const { return exceptionClass_; }
  virtual std::string toString() const;

  // Stack trace support
  void captureStackTrace(Interpreter &interpreter);
  std::vector<std::string> getStackTrace() const { return stackTrace_; }

  // Signal (throw) this exception
  virtual void signal();

protected:
  std::string message_;
  std::string exceptionClass_;
  std::vector<std::string> stackTrace_;
};

/**
 * Zero division error - thrown when dividing by zero
 */
class ZeroDivisionError : public SmalltalkException {
public:
  ZeroDivisionError(const std::string &message = "Division by zero")
      : SmalltalkException(message, "ZeroDivisionError") {}
};

/**
 * Name error - thrown when accessing undefined variables
 */
class NameError : public SmalltalkException {
public:
  NameError(const std::string &variableName)
      : SmalltalkException("Undefined variable: " + variableName, "NameError") {
  }
};

/**
 * Index error - thrown when accessing invalid array/string indices
 */
class IndexError : public SmalltalkException {
public:
  IndexError(const std::string &message = "Index out of bounds")
      : SmalltalkException(message, "IndexError") {}
};

/**
 * Message not understood - thrown when calling unknown methods
 */
class MessageNotUnderstood : public SmalltalkException {
public:
  MessageNotUnderstood(const std::string &receiver, const std::string &selector)
      : SmalltalkException(receiver + " does not understand " + selector,
                           "MessageNotUnderstood") {}
};

/**
 * Argument error - thrown when passing invalid arguments
 */
class ArgumentError : public SmalltalkException {
public:
  ArgumentError(const std::string &message = "Invalid argument")
      : SmalltalkException(message, "ArgumentError") {}
};

/**
 * Runtime error - general runtime exceptions
 */
class RuntimeError : public SmalltalkException {
public:
  RuntimeError(const std::string &message = "Runtime error")
      : SmalltalkException(message, "RuntimeError") {}
};

/**
 * Exception handler for managing Smalltalk exception semantics
 */
class ExceptionHandler {
public:
  // Convert C++ exceptions to Smalltalk exceptions
  static std::unique_ptr<SmalltalkException>
  fromStdException(const std::exception &e);

  // Handle exception during bytecode execution
  static TaggedValue handleException(SmalltalkException &exception,
                                     Interpreter &interpreter);

  // Throw a Smalltalk exception (converts to C++ exception)
  [[noreturn]] static void
  throwException(std::unique_ptr<SmalltalkException> exception);

  // Check if an exception should be caught by a handler
  static bool shouldCatch(const SmalltalkException &exception,
                          const std::string &handlerClass);
};

} // namespace smalltalk