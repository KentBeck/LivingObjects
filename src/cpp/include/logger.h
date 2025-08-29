#pragma once

#include <chrono>
#include <fstream>
#include <iomanip>
#include <iostream>
#include <memory>
#include <sstream>
#include <string>

namespace smalltalk {

enum class LogLevel {
  DEBUG_LEVEL = 0,
  INFO = 1,
  WARN = 2,
  ERROR = 3,
  FATAL = 4
};

class Logger {
public:
  static Logger &getInstance() {
    static Logger instance;
    return instance;
  }

  void setLevel(LogLevel level) { currentLevel = level; }

  LogLevel getLevel() const { return currentLevel; }

  void setOutput(const std::string &filename) {
    logFile = std::make_unique<std::ofstream>(filename, std::ios::app);
    useFile = true;
  }

  void setConsoleOutput(bool enabled) { useConsole = enabled; }

  void log(LogLevel level, const std::string &message,
           const std::string &context = "") {
    if (level < currentLevel)
      return;

    std::string levelStr = getLevelString(level);
    std::string timestamp = getCurrentTimestamp();
    std::string formattedMessage = timestamp + " [" + levelStr + "]";

    if (!context.empty()) {
      formattedMessage += " (" + context + ")";
    }

    formattedMessage += ": " + message;

    if (useConsole) {
      std::cout << formattedMessage << std::endl;
    }

    if (useFile && logFile) {
      *logFile << formattedMessage << std::endl;
      logFile->flush();
    }
  }

  // Convenience methods
  void debug(const std::string &message, const std::string &context = "") {
    log(LogLevel::DEBUG_LEVEL, message, context);
  }

  void info(const std::string &message, const std::string &context = "") {
    log(LogLevel::INFO, message, context);
  }

  void warn(const std::string &message, const std::string &context = "") {
    log(LogLevel::WARN, message, context);
  }

  void error(const std::string &message, const std::string &context = "") {
    log(LogLevel::ERROR, message, context);
  }

  void fatal(const std::string &message, const std::string &context = "") {
    log(LogLevel::FATAL, message, context);
  }

private:
  Logger() : currentLevel(LogLevel::INFO), useConsole(true), useFile(false) {}

  std::string getLevelString(LogLevel level) {
    switch (level) {
    case LogLevel::DEBUG_LEVEL:
      return "DEBUG";
    case LogLevel::INFO:
      return "INFO";
    case LogLevel::WARN:
      return "WARN";
    case LogLevel::ERROR:
      return "ERROR";
    case LogLevel::FATAL:
      return "FATAL";
    default:
      return "UNKNOWN";
    }
  }

  std::string getCurrentTimestamp() {
    auto now = std::chrono::system_clock::now();
    auto time_t = std::chrono::system_clock::to_time_t(now);
    auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(
                  now.time_since_epoch()) %
              1000;

    std::stringstream ss;
    ss << std::put_time(std::localtime(&time_t), "%Y-%m-%d %H:%M:%S");
    ss << '.' << std::setfill('0') << std::setw(3) << ms.count();
    return ss.str();
  }

  LogLevel currentLevel;
  bool useConsole;
  bool useFile;
  std::unique_ptr<std::ofstream> logFile;
};

// Convenient macros for logging
#define LOG_DEBUG(msg, ...) Logger::getInstance().debug(msg, ##__VA_ARGS__)
#define LOG_INFO(msg, ...) Logger::getInstance().info(msg, ##__VA_ARGS__)
#define LOG_WARN(msg, ...) Logger::getInstance().warn(msg, ##__VA_ARGS__)
#define LOG_ERROR(msg, ...) Logger::getInstance().error(msg, ##__VA_ARGS__)
#define LOG_FATAL(msg, ...) Logger::getInstance().fatal(msg, ##__VA_ARGS__)

// VM-specific logging context macros
#define LOG_VM_DEBUG(msg) LOG_DEBUG(msg, "VM")
#define LOG_VM_INFO(msg) LOG_INFO(msg, "VM")
#define LOG_VM_WARN(msg) LOG_WARN(msg, "VM")
#define LOG_VM_ERROR(msg) LOG_ERROR(msg, "VM")

#define LOG_BYTECODE_DEBUG(msg) LOG_DEBUG(msg, "BYTECODE")
#define LOG_MEMORY_DEBUG(msg) LOG_DEBUG(msg, "MEMORY")
#define LOG_GC_DEBUG(msg) LOG_DEBUG(msg, "GC")

} // namespace smalltalk