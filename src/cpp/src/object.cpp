#include "object.h"
#include "smalltalk_class.h"

#include <sstream>

namespace smalltalk {

std::string Object::toString() const {
  std::ostringstream oss;
  if (class_ != nullptr) {
    oss << "a " << class_->getName();
  } else {
    oss << "an Object";
  }
  oss << "@" << std::hex << reinterpret_cast<uintptr_t>(this);
  return oss.str();
}

} // namespace smalltalk