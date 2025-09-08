# Tengri-Lang: A High-Performance Programming Language Concept

[![Status](https://img.shields.io/badge/status-in_development-orange.svg)](https://github.com/DauletBai/tengri-lang)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

**Tengri-Lang** is an experimental programming language designed to test a fundamental hypothesis: can a language architecture derived from the principles of natural agglutinative languages (like Turkic languages) achieve a significant performance increase over traditional computing paradigms?

This project explores a novel approach to compiler design, focusing on minimizing syntactic complexity and eliminating redundant layers of abstraction to boost execution speed.

## 核心理念 (The Core Idea)

The central thesis is that modern programming languages inherit syntactic and structural complexities that create performance overhead. Tengri-Lang is designed around three core principles to combat this:

1.  **Minimalist Agglutinative Syntax:** The grammar is linear and context-dependent, similar to agglutination in linguistics. This drastically simplifies the parsing stage (lexing and AST construction), making the compiler faster and more efficient.
2.  **No Redundant Abstractions:** The language aims to provide a more direct mapping of concepts to execution, removing complex intermediate layers like cumbersome ORMs or verbose protocols. The proposed database and data transfer protocols are integrated at a semantic level.
3.  **Static & Strong Typing from Archetypes:** The type system is built on a small set of universal "archetypes" (e.g., Number, Text, Collection), allowing for aggressive compile-time optimizations.

The primary goal of this research is to achieve a **3x-5x performance improvement** in specific computing domains (like data processing, and high-throughput services) compared to mainstream languages.

## 📂 Repository Structure

This repository documents the entire research and development process.

* **/01_philosophy/**: Contains the original manuscript [**"Tartarus & I" (HTML, Russian)**](./01_philosophy/Tartarus_&_I.html) that outlines the linguistic and philosophical inspiration behind the project.
* **/02_prototype_python/**: A reference implementation of an interpreter written in Python. Used for rapid prototyping of language features.
* **/03_compiler_go/**: The main high-performance implementation of a Tengri-Lang tree-walking interpreter, written in Go. This is the version intended for benchmarking and future development into a full compiler.
* **/04_benchmarks/**: (Coming Soon) A dedicated directory for performance tests and comparison results.

## 🚀 Getting Started

### 1. Python Prototype

Requires Python 3.

```bash
cd 02_prototype_python
python3 main.py

### 2. Go Interpreter

Requires Go 1.24+

```bash
cd 03_compiler_go
go run .

The entry point main.go contains a sample program to execute.