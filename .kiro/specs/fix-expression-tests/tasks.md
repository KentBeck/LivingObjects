# Implementation Plan

- [x] 1. Fix parser keyword message handling for array creation

  - Modify `parseKeywordMessage()` in `simple_parser.cpp` to properly handle `Array new: 3` syntax
  - Ensure the parser doesn't treat the space after `new:` as end of input
  - Add test case to verify `Array new: 3` parses correctly
  - _Requirements: 1.2_

- [ ] 2. Fix parser binary operator handling for string concatenation

  - Verify `parseBinarySelector()` in `simple_parser.cpp` correctly handles the `,` operator
  - Ensure `parseBinaryMessage()` properly chains binary operations
  - Test that `'hello' , ' world'` parses without "Unexpected characters" error
  - _Requirements: 1.3, 5.3_

- [ ] 3. Fix temporary variable declaration parsing

  - Debug `parseTemporaryVariables()` and `isTemporaryVariableDeclaration()` methods
  - Ensure `| x | x := 42. x` syntax is properly recognized and parsed
  - Fix any issues with pipe character handling in the parser
  - _Requirements: 1.4_

- [ ] 4. Fix block parameter parsing

  - Modify `parseBlock()` method to correctly handle `:x |` parameter syntax
  - Ensure block parameters are properly distinguished from block body
  - Test that `[:x | x + 1]` parses correctly without "Unexpected character" error
  - _Requirements: 2.3_

- [ ] 5. Verify and fix string primitive method registration

  - Check that string `size` and `,` methods are properly registered in the class system
  - Debug why `'hello' size` fails with "Method not found: size"
  - Ensure `StringPrimitives::size` and `StringPrimitives::concatenate` are accessible
  - _Requirements: 5.1, 5.2_

- [ ] 6. Fix block execution context management

  - Debug the "Stack pointer below stack start" error in block execution
  - Fix `handleExecuteBlock()` method in interpreter to properly manage stack bounds
  - Ensure block contexts are created with correct stack setup
  - Test that `[3 + 4] value` executes correctly
  - _Requirements: 2.1, 2.2_

- [ ] 7. Implement array creation primitive methods

  - Add `new:` method to Array class that creates arrays with specified size
  - Implement proper array object allocation in the memory manager
  - Ensure `Array new: 3` creates an array with size 3
  - _Requirements: 1.2_

- [ ] 8. Implement array access primitive methods

  - Add `at:` method to Array class for element access (1-based indexing)
  - Add `size` method to Array class to return array length
  - Implement bounds checking and proper error handling
  - Test that `#(1 2 3) at: 2` returns 2 and `#(1 2 3) size` returns 3
  - _Requirements: 3.2, 3.3_

- [ ] 9. Implement boolean conditional methods

  - Add `ifTrue:` method to True and False classes
  - Add `ifFalse:` method to True and False classes
  - Add `ifTrue:ifFalse:` method to True and False classes
  - Implement proper block execution for conditional branches
  - Test that `true ifTrue: [42]` returns 42
  - _Requirements: 4.1, 4.2, 4.3, 4.4_

- [ ] 10. Fix array literal parsing

  - Ensure `parseArrayLiteral()` correctly handles `#(1 2 3)` syntax
  - Verify that array elements are properly parsed and stored
  - Test that array literals create proper array objects
  - _Requirements: 3.1_

- [ ] 11. Verify symbol literal parsing and creation

  - Test that `parseSymbol()` correctly handles `#abc` syntax
  - Ensure symbols are properly interned and created
  - Verify symbol toString() method returns expected format
  - _Requirements: 3.4_

- [ ] 12. Fix parameterized block execution

  - Implement proper argument passing for `[:x | x + 1] value: 5`
  - Ensure block parameters are correctly bound to arguments
  - Fix context creation for parameterized blocks
  - Test that block parameters work correctly in block body
  - _Requirements: 2.4_

- [ ] 13. Run comprehensive expression test verification

  - Execute the full expression test suite
  - Verify all tests marked `shouldPass: true` now pass
  - Document any remaining failures and their root causes
  - Ensure previously passing tests still pass
  - _Requirements: 1.1_

- [ ] 14. Add error handling improvements
  - Improve error messages for parse failures
  - Add proper exception handling for runtime errors
  - Ensure test framework reports clear failure reasons
  - Test edge cases and boundary conditions
  - _Requirements: 1.1_
