# Smalltalk Objects Migration – Progress Log

Purpose: Track the ongoing conversion from C++ container-backed runtime state to pure Smalltalk objects and prepare for a single, generic image serializer. This compresses the working notes and conversation into a clear status snapshot with next actions.

## Summary
- Direction: Represent everything possible as a Smalltalk object (MethodContext, BlockContext, CompiledMethod, Dictionaries, Class metadata). Serialize the heap generically (pointer objects, indexable pointer objects, byte-indexable objects) with a single root (Smalltalk). Execution stack will be stored in Smalltalk (ActiveContext) before image save.
- Today’s aim: Finish MethodDictionary migration and continue removing C++ mirrors while keeping the test suite green.

## What’s Implemented
- Smalltalk globals
  - Global `Smalltalk` Dictionary singleton with primitives and C++ helpers (Globals::get/set). Class registration inserts Classes into `Smalltalk`.
- Dictionary is a real Smalltalk object
  - Storage moved from a C++ map to two Array ivars (keys, values). Primitives `at:`, `at:put:`, `keys`, `size` operate on Arrays. Mirrors Globals when updating `Smalltalk`.
  - Files: `src/cpp/src/primitives/dictionary.cpp`, `src/cpp/src/smalltalk_class.cpp` (Dictionary class size = 2 ivars).
- Blocks and lexical scoping
  - Nested blocks read/write outer temps via name resolution through home chains; shadowing handled in compiler (innermost wins). Extensive tests added; all pass.
- CompiledMethod mirrors
  - Lazy Smalltalk mirrors: `bytecodesBytes` (ByteArray), `literalsArray` (Array with boxed literals), `tempNamesArray` (Array of Symbols). Built on demand by `ensureSmalltalkBacking`.
  - Files: `src/cpp/include/compiled_method.h`, `src/cpp/src/compiled_method.cpp`.
- Class metadata mirrors
  - `ensureSmalltalkMetadata(MemoryManager&)` builds Arrays of Symbols for instance/class var names.
  - Files: `src/cpp/include/smalltalk_class.h`, `src/cpp/src/smalltalk_class.cpp`.
- MethodDictionary Smalltalk mirror (new)
  - Each Class has `methodDictObject_` (Dictionary). `ensureSmalltalkMethodDictionary(MemoryManager&)` bulk-populates keys/values from the current C++ MethodDictionary on demand.
  - Interpreter currently calls `ensureSmalltalkMethodDictionary(memoryManager)` on sends to keep the mirror fresh for serialization/inspection.
  - Files: `src/cpp/include/smalltalk_class.h`, `src/cpp/src/smalltalk_class.cpp`, `src/cpp/src/interpreter.cpp`.
- Tests
  - `run_all_tests.sh` and the expression suite (83/83) are green with the above changes.

## In Progress – MethodDictionary Migration (Phase 1 & 2)
- Phase 1 (write-through on add/remove):
  - Current: `Class::addMethod` and `removeMethod` still update the C++ map and invalidate the Smalltalk mirror for lazy rebuild.
  - Next: update Arrays in `methodDictObject_` in place (or append by reallocating Arrays of size+1) during add/remove so the mirror is always current (dual-write) and avoid invalidation.
- Phase 2 (lookups prefer Smalltalk):
  - Current: `Interpreter::sendMessage` first tries lookup in the Smalltalk mirror (keys/values Arrays), then falls back to the C++ MethodDictionary.
  - `Class::hasMethod` prefers the Smalltalk mirror if present; otherwise falls back to C++.
  - Next: Move `Class::lookupMethod` to prefer Smalltalk mirror as well, then retire C++ lookups once dual-write is in place.

## Next Steps
1. Implement eager write-through to Smalltalk method dictionary in `Class::addMethod/removeMethod`:
   - Ensure `methodDictObject_` exists; ensure keys/values Arrays exist.
   - For add: update value if key exists, else append by allocate-copy Arrays.
   - For remove: find key index, compact Arrays by allocate-copy arrays of size-1.
   - Keep C++ MethodDictionary in sync until fully removed.
2. Switch `Class::lookupMethod` to prefer Smalltalk mirror; retain C++ fallback temporarily.
3. Add a small helper to DRY Array get/set/append semantics used by Dictionary and Class method mirrors.
4. Reflect `ActiveContext` in `Smalltalk` (Smalltalk at: #ActiveContext put: topContext) once the serializer root is ready.
5. Serializer:
   - Two-pass ID assignment, then write generic records (named slots, indexable slots, byte payload).
   - Loader two-pass reconstruct; reinstall `Smalltalk` and resume from `#ActiveContext`.
6. Remove C++ mirrors:
   - MethodDictionary C++ map; CompiledMethod vectors (once all consumers switched); Class var name vectors (optional to keep duplicates for tooling).

## Open Issues / Notes
- Performance: Dictionary and method dictionary lookups are linear over Arrays today; acceptable for bootstrap and tests. Consider Smalltalk-side hashed structure after correctness is in place.
- Tests: The expression suite is configured to count expected exceptions as success. All 83 pass.
- Build notes: macOS toolchain sometimes produced mixed-arch warnings; `make clean` resolves it. Keep an eye on ranlib fat-archive warnings if the toolchain changes.

## How to Verify
- Build and run `src/cpp/run_all_tests.sh` (should be green end-to-end).
- Run `src/cpp/tests/run_expression_tests.sh` and verify summary is 83/83.
- Inspect mirrors interactively: after loading VM, send `class getMethodDictionaryObject` and `ensureSmalltalkMethodDictionary:` paths if you hook in a quick print, or test via a minimal inspection method.

## Quick Timeline (Compressed)
- Fixed nested blocks temp capture, shadowing, and added tests.
- Added Block value primitive and execution; all block tests green.
- Introduced global `Smalltalk` Dictionary and primitives; globals now go through `Smalltalk`.
- Converted Dictionary storage to Smalltalk Arrays (removed C++ map).
- Added CompiledMethod Smalltalk mirrors; later used for serialization.
- Added Class metadata mirrors; started MethodDictionary Smalltalk mirror with lazy populate; interpreter ensures it on sends.
- Started switching lookups to Smalltalk mirror; tests remain green.
- Next: dual-write add/remove to Smalltalk method dictionary; move all lookups to Smalltalk; drop C++ mirrors.

