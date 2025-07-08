# Block Execution Memory Layout

## 1. Just Before #value Gets Sent to a Block

```
STACK (grows upward)
┌─────────────────────┐
│                     │ ← activeContext->stackPointer
├─────────────────────┤
│ BlockContext@0x2050 │ ← receiver for #value message
│ (TaggedValue ptr)   │   (this is what CREATE_BLOCK pushed)
└─────────────────────┘

HEAP OBJECTS
┌──────────────────────────────────────────────────────────────────┐
│ MethodContext (executing "[3 + 4] value")                        │
├──────────────────────────────────────────────────────────────────┤
│ header:                                                          │
│   type: CONTEXT (5)                                              │
│   contextType: METHOD_CONTEXT (0)                               │
│   hash: 0x12345678 (method hash)                                │
│ self: Object@0x1000                                              │
│ sender: nil                                                      │
│ home: nil                                                        │
│ instructionPointer: 15 (at SEND_MESSAGE for #value)             │
│ stackPointer: 0x2100 ───────────────────────────┐               │
└──────────────────────────────────────────────────┼───────────────┘
                                                   │
┌──────────────────────────────────────────────────▼───────────────┐
│ BlockContext (the block [3 + 4])                                │
├──────────────────────────────────────────────────────────────────┤
│ header:                                                          │
│   type: CONTEXT (5)                                              │
│   contextType: BLOCK_CONTEXT (1)                                │
│   hash: 0 (literal index of block method)                       │
│ self: Object@0x1000 (inherited from home)                       │
│ sender: nil                                                      │
│ home: MethodContext@0x2000 ─────────────────────┐               │
│ instructionPointer: 0                           │               │
│ stackPointer: (uninitialized)                   │               │
└──────────────────────────────────────────────────┼───────────────┘
                                                   │
┌──────────────────────────────────────────────────▼───────────────┐
│ CompiledMethod (parent method containing block)                  │
├──────────────────────────────────────────────────────────────────┤
│ header:                                                          │
│   type: METHOD (7)                                               │
│ bytecodes: [                                                     │
│   CREATE_BLOCK, 0, 0, 0, 0, 0, 0, 0, 0,  // literal 0, 0 params│
│   SEND_MESSAGE, 1, 0, 0, 0, 0, 0, 0, 0,  // #value, 0 args     │
│   RETURN_STACK_TOP                                              │
│ ]                                                                │
│ literals: [                                                      │
│   0: CompiledMethod@0x3000, // the block's compiled method      │
│   1: Symbol@0x4000 (#value)                                     │
│ ]                                                                │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│ CompiledMethod (block method for [3 + 4])                       │
├──────────────────────────────────────────────────────────────────┤
│ header:                                                          │
│   type: METHOD (7)                                               │
│ bytecodes: [                                                     │
│   PUSH_LITERAL, 0, 0, 0, 0,     // push 3                       │
│   PUSH_LITERAL, 1, 0, 0, 0,     // push 4                       │
│   SEND_MESSAGE, 2, 0, 0, 0, 1, 0, 0, 0,  // #+, 1 arg          │
│   RETURN_STACK_TOP                                              │
│ ]                                                                │
│ literals: [                                                      │
│   0: SmallInteger(3),                                            │
│   1: SmallInteger(4),                                            │
│   2: Symbol@0x5000 (#+)                                         │
│ ]                                                                │
│ tempVars: [] // no temporaries                                   │
└──────────────────────────────────────────────────────────────────┘

INTERPRETER STATE
┌──────────────────────────────────────────────────────────────────┐
│ activeContext: MethodContext@0x2000                              │
│ currentMethod: CompiledMethod@0x2500 (parent method)             │
│ About to execute: SEND_MESSAGE #value to BlockContext@0x2050    │
└──────────────────────────────────────────────────────────────────┘
```

## 2. When the Next Bytecode is About to be Interpreted (Inside Block)

```
STACK (in new BlockMethodContext)
┌─────────────────────┐
│ (empty)             │ ← blockMethodContext->stackPointer
└─────────────────────┘

HEAP OBJECTS (new context created)
┌──────────────────────────────────────────────────────────────────┐
│ MethodContext (NEW - executing block method)                     │
├──────────────────────────────────────────────────────────────────┤
│ header:                                                          │
│   type: CONTEXT (5)                                              │
│   contextType: METHOD_CONTEXT (0)                               │
│   hash: 0 (no hash for direct execution)                        │
│ self: Object@0x1000 (from home context)                         │
│ sender: MethodContext@0x2000 (the caller)                       │
│ home: nil (method contexts don't have home)                     │
│ instructionPointer: 0 (about to execute first bytecode)         │
│ stackPointer: 0x6100 ───────────────────────────┐               │
│ [temp vars area - empty for [3+4]]              │               │
└──────────────────────────────────────────────────┼───────────────┘
                                                   │
                                                   ▼
                                           (points to stack area)

INTERPRETER STATE
┌──────────────────────────────────────────────────────────────────┐
│ activeContext: MethodContext@0x6000 (new block method context)   │
│ currentMethod: CompiledMethod@0x3000 (the block's method)        │
│ savedContext: MethodContext@0x2000 (parent - will restore later) │
│ savedMethod: CompiledMethod@0x2500 (parent - will restore later) │
│                                                                  │
│ Next bytecode to execute: PUSH_LITERAL at IP=0                  │
│   Will push SmallInteger(3) onto stack                          │
└──────────────────────────────────────────────────────────────────┘

EXECUTION FLOW
1. Block primitive retrieved CompiledMethod@0x3000 from literal 0
2. Created new MethodContext@0x6000 for block execution
3. Set up context with:
   - self from home context
   - sender = current context
   - empty temp vars (block has no parameters or temps)
   - stack pointer at beginning of stack area
4. Switched interpreter state:
   - activeContext → new block method context
   - currentMethod → block's compiled method
   - Saved previous context/method for restoration
5. Ready to interpret first bytecode: PUSH_LITERAL 0
```

## Key Points

1. **Before #value**: The block exists as a BlockContext object on the stack, containing a reference to its home context and the literal index where its CompiledMethod is stored.

2. **After #value primitive**: A new MethodContext is created for executing the block's bytecode, with proper stack setup and the block's CompiledMethod as the current method.

3. **Context Chain**: The sender chain is maintained (block method context → parent method context), allowing proper returns.

4. **Method Access**: The currentMethod in the interpreter switches from the parent method to the block's method, eliminating the need for hash lookups.

5. **Stack Isolation**: Each context has its own stack area, preventing interference between the block's execution and its parent's stack.