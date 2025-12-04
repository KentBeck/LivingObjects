# Smalltalk Primitives Specification

**Version**: 1.0  
**Date**: 2025-01-05  
**Purpose**: Complete specification of all primitive methods with stack diagrams

## Table of Contents

1. [Overview](#overview)
2. [Integer Arithmetic Primitives (1-11)](#integer-arithmetic-primitives)
3. [Array Primitives (60-62)](#array-primitives)
4. [String Primitives (63-67)](#string-primitives)
5. [Object Primitives (70-75, 111)](#object-primitives)
6. [Boolean Conditional Primitives (154-159)](#boolean-conditional-primitives)
7. [Block Execution Primitives (201-202)](#block-execution-primitives)
8. [Dictionary Primitives (700-703)](#dictionary-primitives)
9. [Exception Primitives (1000-1001)](#exception-primitives)
10. [System Primitives (5000+)](#system-primitives)

---

## 1. Overview

### 1.1 Primitive Calling Convention

**Stack Layout Before Primitive Call**:

```
... → ..., receiver, arg1, arg2, ..., argN
```

**Stack Layout After Primitive Call**:

```
... → ..., result
```

### 1.2 Primitive Failure

When a primitive fails:

1. Throw `PrimitiveFailure` exception
2. VM catches exception and executes Smalltalk fallback code
3. If no fallback code exists, propagate error

### 1.3 Stack Diagram Notation

```
Before: ..., value1, value2
After:  ..., result

Where:
  ... = previous stack contents (unchanged)
  value1, value2 = values consumed by primitive
  result = value produced by primitive
```

---

## 2. Integer Arithmetic Primitives (1-11)

### Primitive 1: SmallInteger>>+ (Addition)

**Selector**: `+`  
**Arguments**: 1 (addend)  
**Returns**: Sum as SmallInteger

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(int)
```

**Operation**:

```
result = receiver + arg
```

**Example**:

```smalltalk
3 + 4  "→ 7"

Stack before: ..., 3, 4
Stack after:  ..., 7
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger
- Result overflows 31-bit signed integer range

**Fallback**: Execute Smalltalk implementation if available

---

### Primitive 2: SmallInteger>>- (Subtraction)

**Selector**: `-`  
**Arguments**: 1 (subtrahend)  
**Returns**: Difference as SmallInteger

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(int)
```

**Operation**:

```
result = receiver - arg
```

**Example**:

```smalltalk
10 - 3  "→ 7"

Stack before: ..., 10, 3
Stack after:  ..., 7
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger
- Result overflows 31-bit signed integer range

---

### Primitive 3: SmallInteger>>< (Less Than)

**Selector**: `<`  
**Arguments**: 1 (comparand)  
**Returns**: true or false

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(boolean)
```

**Operation**:

```
result = (receiver < arg) ? true : false
```

**Example**:

```smalltalk
3 < 5   "→ true"
7 < 2   "→ false"

Stack before: ..., 3, 5
Stack after:  ..., true
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger

---

### Primitive 4: SmallInteger>>> (Greater Than)

**Selector**: `>`  
**Arguments**: 1 (comparand)  
**Returns**: true or false

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(boolean)
```

**Operation**:

```
result = (receiver > arg) ? true : false
```

**Example**:

```smalltalk
7 > 2   "→ true"
3 > 5   "→ false"

Stack before: ..., 7, 2
Stack after:  ..., true
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger

---

### Primitive 5: SmallInteger>><= (Less Than or Equal)

**Selector**: `<=`  
**Arguments**: 1 (comparand)  
**Returns**: true or false

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(boolean)
```

**Operation**:

```
result = (receiver <= arg) ? true : false
```

**Example**:

```smalltalk
3 <= 5   "→ true"
5 <= 5   "→ true"
7 <= 2   "→ false"

Stack before: ..., 5, 5
Stack after:  ..., true
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger

---

### Primitive 6: SmallInteger>>>= (Greater Than or Equal)

**Selector**: `>=`  
**Arguments**: 1 (comparand)  
**Returns**: true or false

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(boolean)
```

**Operation**:

```
result = (receiver >= arg) ? true : false
```

**Example**:

```smalltalk
7 >= 2   "→ true"
5 >= 5   "→ true"
3 >= 5   "→ false"

Stack before: ..., 5, 5
Stack after:  ..., true
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger

---

### Primitive 7: SmallInteger>>= (Equality)

**Selector**: `=`  
**Arguments**: 1 (comparand)  
**Returns**: true or false

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(boolean)
```

**Operation**:

```
result = (receiver == arg) ? true : false
```

**Example**:

```smalltalk
5 = 5   "→ true"
3 = 7   "→ false"

Stack before: ..., 5, 5
Stack after:  ..., true
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger

---

### Primitive 8: SmallInteger>>~= (Inequality)

**Selector**: `~=`  
**Arguments**: 1 (comparand)  
**Returns**: true or false

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(boolean)
```

**Operation**:

```
result = (receiver != arg) ? true : false
```

**Example**:

```smalltalk
3 ~= 7   "→ true"
5 ~= 5   "→ false"

Stack before: ..., 3, 7
Stack after:  ..., true
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger

---

### Primitive 9: SmallInteger>>\* (Multiplication)

**Selector**: `*`  
**Arguments**: 1 (multiplier)  
**Returns**: Product as SmallInteger

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(int)
```

**Operation**:

```
result = receiver * arg
```

**Example**:

```smalltalk
6 * 7  "→ 42"

Stack before: ..., 6, 7
Stack after:  ..., 42
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger
- Result overflows 31-bit signed integer range

---

### Primitive 10: SmallInteger>>/ (Division)

**Selector**: `/`  
**Arguments**: 1 (divisor)  
**Returns**: Quotient as SmallInteger

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(int)
```

**Operation**:

```
result = receiver / arg  (integer division, truncated toward zero)
```

**Example**:

```smalltalk
20 / 4   "→ 5"
7 / 2    "→ 3"
-7 / 2   "→ -3"

Stack before: ..., 20, 4
Stack after:  ..., 5
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger
- Argument is zero (signals ZeroDivisionError)

**Special Behavior**:

- Division by zero throws ZeroDivisionError exception
- Integer division truncates toward zero

---

### Primitive 11: SmallInteger>>// (Modulo)

**Selector**: `//`  
**Arguments**: 1 (divisor)  
**Returns**: Remainder as SmallInteger

**Stack Diagram**:

```
Before: ..., receiver(int), arg(int)
After:  ..., result(int)
```

**Operation**:

```
result = receiver % arg  (modulo operation)
```

**Example**:

```smalltalk
10 // 3   "→ 1"
7 // 2    "→ 1"
-7 // 2   "→ -1"

Stack before: ..., 10, 3
Stack after:  ..., 1
```

**Failure Conditions**:

- Receiver is not SmallInteger
- Argument is not SmallInteger
- Argument is zero (signals ZeroDivisionError)

---

## 3. Array Primitives (60-62)

### Primitive 60: Array>>at:

**Selector**: `at:`  
**Arguments**: 1 (index as SmallInteger, 1-based)  
**Returns**: Element at index

**Stack Diagram**:

```
Before: ..., receiver(array), index(int)
After:  ..., element(any)
```

**Operation**:

```
result = receiver[index - 1]  (convert to 0-based indexing)
```

**Example**:

```smalltalk
#(10 20 30) at: 2  "→ 20"

Stack before: ..., #(10 20 30), 2
Stack after:  ..., 20
```

**Failure Conditions**:

- Receiver is not an Array
- Index is not a SmallInteger
- Index < 1 or index > array size (signals IndexError)

**Notes**:

- Smalltalk uses 1-based indexing
- Returns TaggedValue (can be any type)

---

### Primitive 61: Array>>at:put:

**Selector**: `at:put:`  
**Arguments**: 2 (index as SmallInteger, value as any)  
**Returns**: The stored value

**Stack Diagram**:

```
Before: ..., receiver(array), index(int), value(any)
After:  ..., value(any)
```

**Operation**:

```
receiver[index - 1] = value
result = value
```

**Example**:

```smalltalk
| arr |
arr := Array new: 3.
arr at: 2 put: 42.  "→ 42"

Stack before: ..., arr, 2, 42
Stack after:  ..., 42
```

**Failure Conditions**:

- Receiver is not an Array
- Index is not a SmallInteger
- Index < 1 or index > array size (signals IndexError)
- Array is immutable (signals error)

**Notes**:

- Returns the value that was stored
- Assignment expression returns the assigned value

---

### Primitive 62: Array>>size

**Selector**: `size`  
**Arguments**: 0  
**Returns**: Number of elements as SmallInteger

**Stack Diagram**:

```
Before: ..., receiver(array)
After:  ..., size(int)
```

**Operation**:

```
result = receiver.length
```

**Example**:

```smalltalk
#(1 2 3 4 5) size  "→ 5"

Stack before: ..., #(1 2 3 4 5)
Stack after:  ..., 5
```

**Failure Conditions**:

- Receiver is not an Array

**Notes**:

- Returns SmallInteger
- Empty array returns 0

---

## 4. String Primitives (63-67)

### Primitive 63: String>>at:

**Selector**: `at:`  
**Arguments**: 1 (index as SmallInteger, 1-based)  
**Returns**: Character at index as SmallInteger (Unicode code point)

**Stack Diagram**:

```
Before: ..., receiver(string), index(int)
After:  ..., character(int)
```

**Operation**:

```
result = receiver[index - 1].codePoint  (Unicode code point as integer)
```

**Example**:

```smalltalk
'hello' at: 2  "→ 101 (code point for 'e')"

Stack before: ..., 'hello', 2
Stack after:  ..., 101
```

**Failure Conditions**:

- Receiver is not a String
- Index is not a SmallInteger
- Index < 1 or index > string length (signals IndexError)

**Notes**:

- Returns Unicode code point as SmallInteger
- Smalltalk uses 1-based indexing

---

### Primitive 64: String>>at:put:

**Selector**: `at:put:`  
**Arguments**: 2 (index as SmallInteger, character as SmallInteger)  
**Returns**: The character value

**Stack Diagram**:

```
Before: ..., receiver(string), index(int), char(int)
After:  ..., char(int)
```

**Operation**:

```
receiver[index - 1] = char.toCharacter
result = char
```

**Example**:

```smalltalk
| str |
str := String new: 5.
str at: 1 put: 65.  "→ 65 (sets first char to 'A')"

Stack before: ..., str, 1, 65
Stack after:  ..., 65
```

**Failure Conditions**:

- Receiver is not a String
- Index is not a SmallInteger
- Character is not a SmallInteger
- Index < 1 or index > string length (signals IndexError)
- String is immutable (signals error)
- Character code point is invalid

---

### Primitive 65: String>>, (Concatenation)

**Selector**: `,`  
**Arguments**: 1 (other string)  
**Returns**: New concatenated string

**Stack Diagram**:

```
Before: ..., receiver(string), arg(string)
After:  ..., result(string)
```

**Operation**:

```
result = receiver + arg  (string concatenation)
```

**Example**:

```smalltalk
'hello' , ' world'  "→ 'hello world'"

Stack before: ..., 'hello', ' world'
Stack after:  ..., 'hello world'
```

**Failure Conditions**:

- Receiver is not a String
- Argument is not a String
- Memory allocation fails

**Notes**:

- Creates new String object
- Does not modify receiver or argument

---

### Primitive 66: String>>size

**Selector**: `size`  
**Arguments**: 0  
**Returns**: Number of characters as SmallInteger

**Stack Diagram**:

```
Before: ..., receiver(string)
After:  ..., size(int)
```

**Operation**:

```
result = receiver.length
```

**Example**:

```smalltalk
'hello' size  "→ 5"
'' size       "→ 0"

Stack before: ..., 'hello'
Stack after:  ..., 5
```

**Failure Conditions**:

- Receiver is not a String

**Notes**:

- Returns number of characters (not bytes)
- Empty string returns 0

---

### Primitive 67: String>>asSymbol

**Selector**: `asSymbol`  
**Arguments**: 0  
**Returns**: Interned symbol

**Stack Diagram**:

```
Before: ..., receiver(string)
After:  ..., symbol(symbol)
```

**Operation**:

```
result = Symbol.intern(receiver)
```

**Example**:

```smalltalk
'hello' asSymbol  "→ #hello"
'test' asSymbol == 'test' asSymbol  "→ true (same object)"

Stack before: ..., 'hello'
Stack after:  ..., #hello
```

**Failure Conditions**:

- Receiver is not a String

**Notes**:

- Symbols are interned (unique per string)
- Same string always returns same symbol object
- Symbol equality is pointer equality

---

## 5. Object Primitives (70-75, 111)

### Primitive 70: Object>>new

**Selector**: `new`  
**Arguments**: 0  
**Returns**: New instance of receiver class

**Stack Diagram**:

```
Before: ..., receiver(class)
After:  ..., instance(object)
```

**Operation**:

```
1. Allocate object with class's instance variable count
2. Initialize all instance variables to nil
3. Send #initialize message to new instance
4. Return instance
```

**Example**:

```smalltalk
Object new  "→ new Object instance"

Stack before: ..., Object
Stack after:  ..., anObject
```

**Failure Conditions**:

- Receiver is not a Class
- Memory allocation fails

**Notes**:

- Sends `initialize` message after creation
- Instance variables initialized to nil

---

### Primitive 71: Object>>basicNew

**Selector**: `basicNew`  
**Arguments**: 0  
**Returns**: New instance of receiver class

**Stack Diagram**:

```
Before: ..., receiver(class)
After:  ..., instance(object)
```

**Operation**:

```
1. Allocate object with class's instance variable count
2. Initialize all instance variables to nil
3. Return instance (NO initialize message)
```

**Example**:

```smalltalk
Object basicNew  "→ new Object instance (not initialized)"

Stack before: ..., Object
Stack after:  ..., anObject
```

**Failure Conditions**:

- Receiver is not a Class
- Memory allocation fails

**Notes**:

- Does NOT send `initialize` message
- Used when custom initialization needed
- Faster than `new` for simple objects

---

### Primitive 72: Object>>basicNew: (Variable-sized objects)

**Selector**: `basicNew:`  
**Arguments**: 1 (size as SmallInteger)  
**Returns**: New instance with specified size

**Stack Diagram**:

```
Before: ..., receiver(class), size(int)
After:  ..., instance(object)
```

**Operation**:

```
1. Allocate object with specified size
2. Initialize all slots to nil
3. Return instance (NO initialize message)
```

**Example**:

```smalltalk
Array new: 5    "→ array with 5 elements"
String new: 10  "→ string with 10 characters"

Stack before: ..., Array, 5
Stack after:  ..., anArray
```

**Failure Conditions**:

- Receiver is not a Class
- Size is not a SmallInteger
- Size is negative (signals ArgumentError)
- Memory allocation fails

**Notes**:

- Used for Array, String, ByteArray creation
- Does NOT send `initialize` message
- All elements/characters initialized to nil/zero

---

### Primitive 75: Object>>identityHash

**Selector**: `identityHash`  
**Arguments**: 0  
**Returns**: Identity hash as SmallInteger

**Stack Diagram**:

```
Before: ..., receiver(any)
After:  ..., hash(int)
```

**Operation**:

```
result = receiver.objectHeader.hash
```

**Example**:

```smalltalk
Object new identityHash  "→ 12345 (some hash value)"

Stack before: ..., anObject
Stack after:  ..., 12345
```

**Failure Conditions**:

- None (works for all objects)

**Notes**:

- Returns hash from object header
- Hash is stable for object lifetime
- Used for hashing in dictionaries
- Immediate values (SmallInteger, Boolean, nil) have computed hashes

---

### Primitive 111: Object>>class

**Selector**: `class`  
**Arguments**: 0  
**Returns**: Receiver's class

**Stack Diagram**:

```
Before: ..., receiver(any)
After:  ..., class(class)
```

**Operation**:

```
result = receiver.class
```

**Example**:

```smalltalk
42 class        "→ SmallInteger"
'hello' class   "→ String"
true class      "→ True"

Stack before: ..., 42
Stack after:  ..., SmallInteger
```

**Failure Conditions**:

- None (all objects have a class)

**Notes**:

- Works for immediate values (SmallInteger, Boolean, nil)
- Returns Class object
- Used for type checking and introspection

---

## 6. Boolean Conditional Primitives (154-159)

### Primitive 154: True>>ifTrue:

**Selector**: `ifTrue:`  
**Arguments**: 1 (block)  
**Returns**: Result of evaluating block

**Stack Diagram**:

```
Before: ..., receiver(true), block(block)
After:  ..., result(any)
```

**Operation**:

```
1. Verify receiver is true
2. Execute block
3. Return block's result
```

**Example**:

```smalltalk
true ifTrue: [42]  "→ 42"

Stack before: ..., true, [42]
Stack after:  ..., 42
```

**Failure Conditions**:

- Receiver is not true (exactly)
- Argument is not a Block

**Notes**:

- Block is evaluated
- Returns block's result
- Short-circuit evaluation (block only evaluated if true)

---

### Primitive 155: True>>ifFalse:

**Selector**: `ifFalse:`  
**Arguments**: 1 (block)  
**Returns**: nil

**Stack Diagram**:

```
Before: ..., receiver(true), block(block)
After:  ..., nil
```

**Operation**:

```
1. Verify receiver is true
2. Do NOT execute block
3. Return nil
```

**Example**:

```smalltalk
true ifFalse: [42]  "→ nil"

Stack before: ..., true, [42]
Stack after:  ..., nil
```

**Failure Conditions**:

- Receiver is not true (exactly)
- Argument is not a Block

**Notes**:

- Block is NOT evaluated
- Always returns nil
- Short-circuit evaluation

---

### Primitive 156: True>>ifTrue:ifFalse:

**Selector**: `ifTrue:ifFalse:`  
**Arguments**: 2 (trueBlock, falseBlock)  
**Returns**: Result of evaluating trueBlock

**Stack Diagram**:

```
Before: ..., receiver(true), trueBlock(block), falseBlock(block)
After:  ..., result(any)
```

**Operation**:

```
1. Verify receiver is true
2. Execute trueBlock
3. Return trueBlock's result
```

**Example**:

```smalltalk
true ifTrue: [42] ifFalse: [99]  "→ 42"

Stack before: ..., true, [42], [99]
Stack after:  ..., 42
```

**Failure Conditions**:

- Receiver is not true (exactly)
- Arguments are not Blocks

**Notes**:

- Only trueBlock is evaluated
- falseBlock is ignored
- Short-circuit evaluation

---

### Primitive 157: False>>ifTrue:

**Selector**: `ifTrue:`  
**Arguments**: 1 (block)  
**Returns**: nil

**Stack Diagram**:

```
Before: ..., receiver(false), block(block)
After:  ..., nil
```

**Operation**:

```
1. Verify receiver is false
2. Do NOT execute block
3. Return nil
```

**Example**:

```smalltalk
false ifTrue: [42]  "→ nil"

Stack before: ..., false, [42]
Stack after:  ..., nil
```

**Failure Conditions**:

- Receiver is not false (exactly)
- Argument is not a Block

**Notes**:

- Block is NOT evaluated
- Always returns nil
- Short-circuit evaluation

---

### Primitive 158: False>>ifFalse:

**Selector**: `ifFalse:`  
**Arguments**: 1 (block)  
**Returns**: Result of evaluating block

**Stack Diagram**:

```
Before: ..., receiver(false), block(block)
After:  ..., result(any)
```

**Operation**:

```
1. Verify receiver is false
2. Execute block
3. Return block's result
```

**Example**:

```smalltalk
false ifFalse: [42]  "→ 42"

Stack before: ..., false, [42]
Stack after:  ..., 42
```

**Failure Conditions**:

- Receiver is not false (exactly)
- Argument is not a Block

**Notes**:

- Block is evaluated
- Returns block's result
- Short-circuit evaluation

---

### Primitive 159: False>>ifTrue:ifFalse:

**Selector**: `ifTrue:ifFalse:`  
**Arguments**: 2 (trueBlock, falseBlock)  
**Returns**: Result of evaluating falseBlock

**Stack Diagram**:

```
Before: ..., receiver(false), trueBlock(block), falseBlock(block)
After:  ..., result(any)
```

**Operation**:

```
1. Verify receiver is false
2. Execute falseBlock
3. Return falseBlock's result
```

**Example**:

```smalltalk
false ifTrue: [42] ifFalse: [99]  "→ 99"

Stack before: ..., false, [42], [99]
Stack after:  ..., 99
```

**Failure Conditions**:

- Receiver is not false (exactly)
- Arguments are not Blocks

**Notes**:

- Only falseBlock is evaluated
- trueBlock is ignored
- Short-circuit evaluation

---

## 7. Block Execution Primitives (201-202)

### Primitive 201: Block>>value

**Selector**: `value`  
**Arguments**: 0  
**Returns**: Result of block execution

**Stack Diagram**:

```
Before: ..., receiver(block)
After:  ..., result(any)
```

**Operation**:

```
1. Verify receiver is BlockContext
2. Create execution context for block
3. Set sender to current context
4. Execute block bytecode
5. Return result
```

**Example**:

```smalltalk
[3 + 4] value  "→ 7"

Stack before: ..., [3 + 4]
Stack after:  ..., 7
```

**Failure Conditions**:

- Receiver is not a BlockContext
- Block expects arguments (argument count mismatch)

**Notes**:

- Block executes with home context's self
- Block can access home context variables
- Creates new execution context
- Returns to sender when complete

---

### Primitive 202: Block>>value:

**Selector**: `value:`  
**Arguments**: 1 (argument)  
**Returns**: Result of block execution

**Stack Diagram**:

```
Before: ..., receiver(block), arg(any)
After:  ..., result(any)
```

**Operation**:

```
1. Verify receiver is BlockContext
2. Create execution context for block
3. Copy argument to first temporary variable
4. Set sender to current context
5. Execute block bytecode
6. Return result
```

**Example**:

```smalltalk
[:x | x + 1] value: 5  "→ 6"

Stack before: ..., [:x | x + 1], 5
Stack after:  ..., 6
```

**Failure Conditions**:

- Receiver is not a BlockContext
- Block expects different number of arguments
- Argument count mismatch

**Notes**:

- Argument is bound to block parameter
- Block executes with home context's self
- Can access home context variables
- Returns to sender when complete

---

### Block>>value:value: (Not a primitive, implemented in Smalltalk)

**Selector**: `value:value:`  
**Arguments**: 2 (arg1, arg2)  
**Returns**: Result of block execution

**Stack Diagram**:

```
Before: ..., receiver(block), arg1(any), arg2(any)
After:  ..., result(any)
```

**Implementation**:

```smalltalk
Block>>value: arg1 value: arg2
    "Execute block with two arguments"
    <primitive: 202>  "Uses same primitive as value:"
    ^self valueWithArguments: {arg1. arg2}
```

**Example**:

```smalltalk
[:x :y | x + y] value: 3 value: 4  "→ 7"

Stack before: ..., [:x :y | x + y], 3, 4
Stack after:  ..., 7
```

**Notes**:

- Typically implemented using primitive 202 with multiple arguments
- Or implemented in Smalltalk using valueWithArguments:

---

## 8. Dictionary Primitives (700-703)

### Primitive 700: Dictionary>>at:

**Selector**: `at:`  
**Arguments**: 1 (key)  
**Returns**: Value associated with key

**Stack Diagram**:

```
Before: ..., receiver(dict), key(any)
After:  ..., value(any)
```

**Operation**:

```
result = receiver.lookup(key)
if not found: signal KeyNotFound
```

**Example**:

```smalltalk
| dict |
dict := Dictionary new.
dict at: #name put: 'Alice'.
dict at: #name  "→ 'Alice'"

Stack before: ..., dict, #name
Stack after:  ..., 'Alice'
```

**Failure Conditions**:

- Receiver is not a Dictionary
- Key not found (signals KeyNotFound)

**Notes**:

- Uses key equality (=) for lookup
- Signals exception if key not found
- Use at:ifAbsent: for default values

---

### Primitive 701: Dictionary>>at:put:

**Selector**: `at:put:`  
**Arguments**: 2 (key, value)  
**Returns**: The stored value

**Stack Diagram**:

```
Before: ..., receiver(dict), key(any), value(any)
After:  ..., value(any)
```

**Operation**:

```
receiver.store(key, value)
result = value
```

**Example**:

```smalltalk
| dict |
dict := Dictionary new.
dict at: #name put: 'Alice'  "→ 'Alice'"

Stack before: ..., dict, #name, 'Alice'
Stack after:  ..., 'Alice'
```

**Failure Conditions**:

- Receiver is not a Dictionary
- Dictionary is immutable (signals error)

**Notes**:

- Creates new entry if key doesn't exist
- Updates existing entry if key exists
- Returns the value that was stored

---

### Primitive 702: Dictionary>>keys

**Selector**: `keys`  
**Arguments**: 0  
**Returns**: Collection of all keys

**Stack Diagram**:

```
Before: ..., receiver(dict)
After:  ..., keys(collection)
```

**Operation**:

```
result = receiver.allKeys()
```

**Example**:

```smalltalk
| dict |
dict := Dictionary new.
dict at: #a put: 1.
dict at: #b put: 2.
dict keys  "→ #(#a #b) or similar collection"

Stack before: ..., dict
Stack after:  ..., keys
```

**Failure Conditions**:

- Receiver is not a Dictionary

**Notes**:

- Returns collection (typically Array or OrderedCollection)
- Order of keys is implementation-defined
- Empty dictionary returns empty collection

---

### Primitive 703: Dictionary>>size

**Selector**: `size`  
**Arguments**: 0  
**Returns**: Number of key-value pairs as SmallInteger

**Stack Diagram**:

```
Before: ..., receiver(dict)
After:  ..., size(int)
```

**Operation**:

```
result = receiver.entryCount()
```

**Example**:

```smalltalk
| dict |
dict := Dictionary new.
dict at: #a put: 1.
dict at: #b put: 2.
dict size  "→ 2"

Stack before: ..., dict
Stack after:  ..., 2
```

**Failure Conditions**:

- Receiver is not a Dictionary

**Notes**:

- Returns SmallInteger
- Empty dictionary returns 0

---

## 9. Exception Primitives (1000-1001)

### Primitive 1000: Exception>>mark (Handler Marker)

**Selector**: N/A (internal use)  
**Arguments**: 0  
**Returns**: Never returns (always fails)

**Stack Diagram**:

```
Before: ..., receiver(any)
After:  (primitive always fails)
```

**Operation**:

```
1. Mark current context as exception handler
2. Always fail (fall back to Smalltalk code)
```

**Example**:

```smalltalk
"Internal use in exception handling implementation"
[
    "protected code"
] on: Error do: [:ex |
    "handler code"
]
```

**Failure Conditions**:

- Always fails (by design)

**Notes**:

- Used internally by exception handling mechanism
- Marks context for exception handler search
- Never executes successfully
- Fallback code implements actual exception handling

---

### Primitive 1001: Exception>>signal

**Selector**: `signal`  
**Arguments**: 0  
**Returns**: Never returns normally (throws exception)

**Stack Diagram**:

```
Before: ..., receiver(exception)
After:  (exception thrown, stack unwound)
```

**Operation**:

```
1. Create exception object
2. Search context chain for handler
3. If found, unwind to handler context
4. Execute handler block with exception
5. If not found, terminate with unhandled exception
```

**Example**:

```smalltalk
Error signal: 'Something went wrong'

Stack before: ..., anError
Stack after:  (exception thrown)
```

**Failure Conditions**:

- Receiver is not an Exception

**Notes**:

- Throws exception and unwinds stack
- Searches for matching exception handler
- Executes ensure: blocks during unwinding
- If no handler found, terminates program

---

### Exception Handling Stack Behavior

**Normal Execution**:

```
Before: ..., context1, context2, context3
After:  ..., context1, context2, context3, result
```

**Exception Thrown**:

```
Before: ..., context1(handler), context2, context3(throws)
After:  ..., context1(handler), handlerResult

(context2 and context3 unwound)
```

**With ensure: blocks**:

```
Before: ..., context1(handler), context2(ensure), context3(throws)
After:  ..., context1(handler), handlerResult

(ensure block in context2 executed during unwind)
```

---

## 10. System Primitives (5000+)

### Primitive 5000: SystemLoader>>start:

**Selector**: `start:`  
**Arguments**: 1 (source directory path)  
**Returns**: Result of bootstrap process

**Stack Diagram**:

```
Before: ..., receiver(systemLoader), path(string)
After:  ..., result(any)
```

**Operation**:

```
1. Load Smalltalk source files from directory
2. Compile classes and methods
3. Build initial image
4. Return result
```

**Example**:

```smalltalk
SystemLoader new start: 'src/'

Stack before: ..., aSystemLoader, 'src/'
Stack after:  ..., result
```

**Failure Conditions**:

- Receiver is not a SystemLoader
- Path is not a String
- Directory not found
- Compilation errors

**Notes**:

- Used for bootstrap process
- Loads and compiles Smalltalk source
- Implementation-specific

---

### Primitive 5100: Compiler>>compile:in:

**Selector**: `compile:in:`  
**Arguments**: 2 (source code string, target class)  
**Returns**: CompiledMethod

**Stack Diagram**:

```
Before: ..., receiver(compiler), source(string), class(class)
After:  ..., compiledMethod(method)
```

**Operation**:

```
1. Parse source code
2. Generate bytecode
3. Create CompiledMethod
4. Return compiled method
```

**Example**:

```smalltalk
Compiler new compile: 'factorial: n
    n <= 1 ifTrue: [^1].
    ^n * (self factorial: n - 1)' in: Integer

Stack before: ..., aCompiler, sourceString, Integer
Stack after:  ..., aCompiledMethod
```

**Failure Conditions**:

- Receiver is not a Compiler
- Source is not a String
- Target is not a Class
- Syntax errors in source
- Compilation errors

**Notes**:

- Used for dynamic compilation
- Returns CompiledMethod object
- Can be installed in class

---

## 11. Primitive Summary Table

### Complete Primitive List

| Number | Selector          | Receiver     | Args           | Returns    | Description               |
| ------ | ----------------- | ------------ | -------------- | ---------- | ------------------------- |
| 1      | `+`               | SmallInteger | 1 int          | int        | Addition                  |
| 2      | `-`               | SmallInteger | 1 int          | int        | Subtraction               |
| 3      | `<`               | SmallInteger | 1 int          | bool       | Less than                 |
| 4      | `>`               | SmallInteger | 1 int          | bool       | Greater than              |
| 5      | `<=`              | SmallInteger | 1 int          | bool       | Less than or equal        |
| 6      | `>=`              | SmallInteger | 1 int          | bool       | Greater than or equal     |
| 7      | `=`               | SmallInteger | 1 int          | bool       | Equality                  |
| 8      | `~=`              | SmallInteger | 1 int          | bool       | Inequality                |
| 9      | `*`               | SmallInteger | 1 int          | int        | Multiplication            |
| 10     | `/`               | SmallInteger | 1 int          | int        | Division                  |
| 11     | `//`              | SmallInteger | 1 int          | int        | Modulo                    |
| 60     | `at:`             | Array        | 1 int          | any        | Element access            |
| 61     | `at:put:`         | Array        | 2 int,any      | any        | Element store             |
| 62     | `size`            | Array        | 0              | int        | Array size                |
| 63     | `at:`             | String       | 1 int          | int        | Character access          |
| 64     | `at:put:`         | String       | 2 int,int      | int        | Character store           |
| 65     | `,`               | String       | 1 string       | string     | Concatenation             |
| 66     | `size`            | String       | 0              | int        | String length             |
| 67     | `asSymbol`        | String       | 0              | symbol     | Convert to symbol         |
| 70     | `new`             | Class        | 0              | object     | Create instance           |
| 71     | `basicNew`        | Class        | 0              | object     | Create instance (no init) |
| 72     | `basicNew:`       | Class        | 1 int          | object     | Create sized instance     |
| 75     | `identityHash`    | Object       | 0              | int        | Identity hash             |
| 111    | `class`           | Object       | 0              | class      | Get object's class        |
| 154    | `ifTrue:`         | True         | 1 block        | any        | Execute if true           |
| 155    | `ifFalse:`        | True         | 1 block        | nil        | Don't execute             |
| 156    | `ifTrue:ifFalse:` | True         | 2 block,block  | any        | Execute true block        |
| 157    | `ifTrue:`         | False        | 1 block        | nil        | Don't execute             |
| 158    | `ifFalse:`        | False        | 1 block        | any        | Execute if false          |
| 159    | `ifTrue:ifFalse:` | False        | 2 block,block  | any        | Execute false block       |
| 201    | `value`           | Block        | 0              | any        | Execute block             |
| 202    | `value:`          | Block        | 1 any          | any        | Execute block with arg    |
| 700    | `at:`             | Dictionary   | 1 any          | any        | Lookup value              |
| 701    | `at:put:`         | Dictionary   | 2 any,any      | any        | Store value               |
| 702    | `keys`            | Dictionary   | 0              | collection | Get all keys              |
| 703    | `size`            | Dictionary   | 0              | int        | Entry count               |
| 1000   | (internal)        | Exception    | 0              | -          | Handler marker            |
| 1001   | `signal`          | Exception    | 0              | -          | Throw exception           |
| 5000   | `start:`          | SystemLoader | 1 string       | any        | Bootstrap system          |
| 5100   | `compile:in:`     | Compiler     | 2 string,class | method     | Compile method            |

---

## 12. Common Primitive Patterns

### Pattern 1: Binary Arithmetic Operation

```
Stack before: ..., receiver(int), arg(int)
Stack after:  ..., result(int)

Operation: result = receiver OP arg
Failures: Non-integer receiver/arg, overflow
```

**Primitives**: 1, 2, 9, 10, 11

---

### Pattern 2: Binary Comparison Operation

```
Stack before: ..., receiver(int), arg(int)
Stack after:  ..., result(bool)

Operation: result = (receiver CMP arg) ? true : false
Failures: Non-integer receiver/arg
```

**Primitives**: 3, 4, 5, 6, 7, 8

---

### Pattern 3: Indexed Access

```
Stack before: ..., receiver(collection), index(int)
Stack after:  ..., element(any)

Operation: result = receiver[index - 1]
Failures: Non-collection receiver, non-integer index, out of bounds
```

**Primitives**: 60, 63

---

### Pattern 4: Indexed Store

```
Stack before: ..., receiver(collection), index(int), value(any)
Stack after:  ..., value(any)

Operation: receiver[index - 1] = value; result = value
Failures: Non-collection receiver, non-integer index, out of bounds, immutable
```

**Primitives**: 61, 64

---

### Pattern 5: Size Query

```
Stack before: ..., receiver(collection)
Stack after:  ..., size(int)

Operation: result = receiver.length
Failures: Non-collection receiver
```

**Primitives**: 62, 66, 703

---

### Pattern 6: Object Creation

```
Stack before: ..., receiver(class), [size(int)]
Stack after:  ..., instance(object)

Operation: Allocate and initialize new instance
Failures: Non-class receiver, allocation failure
```

**Primitives**: 70, 71, 72

---

### Pattern 7: Conditional Execution

```
Stack before: ..., receiver(bool), block(block), [block(block)]
Stack after:  ..., result(any)

Operation: Execute appropriate block based on receiver
Failures: Non-boolean receiver, non-block argument
```

**Primitives**: 154-159

---

### Pattern 8: Block Execution

```
Stack before: ..., receiver(block), [args...]
Stack after:  ..., result(any)

Operation: Execute block with arguments
Failures: Non-block receiver, argument count mismatch
```

**Primitives**: 201, 202

---

## 13. Implementation Notes

### 13.1 Primitive Failure Handling

When a primitive fails:

1. **Catch PrimitiveFailure exception**
2. **Check for Smalltalk fallback code** in method
3. **If fallback exists**: Execute Smalltalk bytecode
4. **If no fallback**: Propagate error to caller

### 13.2 Stack Management

Primitives must:

- **Pop arguments** from stack (including receiver)
- **Push result** onto stack
- **Maintain stack balance** (pop N+1, push 1)
- **Handle errors** without corrupting stack

### 13.3 Type Checking

Primitives should:

- **Verify receiver type** before operation
- **Verify argument types** before operation
- **Fail gracefully** with PrimitiveFailure
- **Provide clear error messages**

### 13.4 Memory Management

Primitives that allocate must:

- **Use memory manager** for allocation
- **Handle allocation failures** gracefully
- **Initialize allocated objects** properly
- **Return properly tagged values**

### 13.5 Exception Safety

Primitives should:

- **Not leak resources** on failure
- **Maintain VM consistency** on error
- **Clean up partial state** before failing
- **Use RAII** for resource management

---

## 14. Testing Primitives

### 14.1 Test Coverage

Each primitive should have tests for:

- **Normal operation** with valid inputs
- **Boundary conditions** (min/max values, empty collections)
- **Error conditions** (wrong types, out of bounds)
- **Edge cases** (overflow, division by zero, nil)

### 14.2 Example Test Cases

**Integer Addition (Primitive 1)**:

```smalltalk
testIntegerAddition
    self assert: 3 + 4 equals: 7.
    self assert: 0 + 0 equals: 0.
    self assert: -5 + 3 equals: -2.
    self should: [3 + 'string'] raise: PrimitiveFailure.
```

**Array Access (Primitive 60)**:

```smalltalk
testArrayAccess
    | arr |
    arr := #(10 20 30).
    self assert: (arr at: 1) equals: 10.
    self assert: (arr at: 3) equals: 30.
    self should: [arr at: 0] raise: IndexError.
    self should: [arr at: 4] raise: IndexError.
```

**Boolean Conditionals (Primitives 154-159)**:

```smalltalk
testBooleanConditionals
    self assert: (true ifTrue: [42]) equals: 42.
    self assert: (true ifFalse: [42]) equals: nil.
    self assert: (false ifTrue: [42]) equals: nil.
    self assert: (false ifFalse: [42]) equals: 42.
    self assert: (true ifTrue: [1] ifFalse: [2]) equals: 1.
    self assert: (false ifTrue: [1] ifFalse: [2]) equals: 2.
```

---

## 15. Appendix: Quick Reference

### Primitive Number Ranges

- **1-11**: Integer operations
- **60-67**: Collection operations (Array, String)
- **70-75, 111**: Object operations
- **154-159**: Boolean conditionals
- **201-202**: Block execution
- **700-703**: Dictionary operations
- **1000-1001**: Exception handling
- **5000+**: System/bootstrap operations

### Common Failure Reasons

1. **Type mismatch**: Receiver or argument wrong type
2. **Out of bounds**: Index < 1 or > size
3. **Division by zero**: Divisor is 0
4. **Overflow**: Result exceeds 31-bit range
5. **Not found**: Key not in dictionary
6. **Immutable**: Attempt to modify immutable object
7. **Allocation failure**: Out of memory
8. **Argument count**: Wrong number of arguments

### Stack Effect Summary

- **Unary**: `..., receiver → ..., result` (pop 1, push 1)
- **Binary**: `..., receiver, arg → ..., result` (pop 2, push 1)
- **Ternary**: `..., receiver, arg1, arg2 → ..., result` (pop 3, push 1)

---

**End of Primitives Specification**
