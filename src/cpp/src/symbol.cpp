#include "symbol.h"

namespace smalltalk {

// Static member definition
std::unordered_map<std::string, std::unique_ptr<Symbol>> Symbol::symbolTable_;

Symbol *Symbol::intern(const std::string &name) {
  auto it = symbolTable_.find(name);
  if (it != symbolTable_.end()) {
    return it->second.get();
  }

  // Create new symbol and add to table
  auto symbol = std::unique_ptr<Symbol>(new Symbol(name));
  Symbol *symbolPtr = symbol.get();
  symbolTable_[name] = std::move(symbol);

  return symbolPtr;
}

} // namespace smalltalk