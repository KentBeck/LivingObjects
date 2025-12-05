# Base Image Canonical Smalltalk Code

This directory contains the canonical Smalltalk source code that is loaded to create the base image. The code is organized by package/category.

## Directory Structure

- `Kernel-Boolean/` - Boolean classes (True, False)
- `Kernel-Objects/` - Core object classes (Object, Class, etc.)
- `Kernel-Numbers/` - Number classes (SmallInteger, Float, etc.)
- `Collections-Sequenceable/` - Array, String, OrderedCollection
- `Collections-Unordered/` - Set, Dictionary
- `Kernel-Methods/` - CompiledMethod, BlockClosure
- `Kernel-Exceptions/` - Exception handling classes

## Loading Order

Classes should be loaded in dependency order:

1. Object (root of hierarchy)
2. Boolean classes (True, False)
3. Number classes
4. Collections
5. Methods and blocks
6. Exceptions

## File Format

Each `.st` file contains methods for a single class. Methods are separated by blank lines and follow standard Smalltalk syntax.
