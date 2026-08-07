# Goduck: NestJS-Inspired Backend DSL for Go

Goduck is a backend Domain Specific Language (DSL) that allows developers to write modular, TypeScript-inspired backend applications that compile into highly optimized Go servers.

## Features

- **TypeScript-Inspired Syntax**: Write code using modular components, classes, and decorators.
- **Recursive Multi-File Compiler**: Combines files in your project directory into a unified Abstract Syntax Tree (AST).
- **Strong Type Checking**: Prevents compiling errors by statically analyzing types before generating Go code.
- **Auto-wired Dependency Injection**: Automatically configures dependencies of controllers and services at compile time.
- **Auto-generated Gin Server**: Automatically parses endpoints, maps parameters, and binds JSON requests.

---

## Getting Started

### 1. Compile the Goduck Compiler

To compile the `goduck` compiler on your machine:

```bash
go build -o goduck.exe ./cmd/goduck
```

### 2. Scaffold a New Project

Initialize a fully structured Goduck project:

```bash
./goduck.exe new myapp
```

This generates the workspace layout:

```
myapp/
├── .env.example
├── .gitignore
├── README.md
└── src/
    ├── app.duck
    ├── main.duck
    ├── auth/
    │   ├── auth.controller.duck
    │   └── auth.service.duck
    └── users/
        ├── user.duck
        ├── users.dto.duck
        ├── users.service.duck
        └── users.controller.duck
```

### 3. Build & Run the App

To compile the newly created project:

```bash
cd myapp
../goduck.exe build src
```

Run the compiled Go backend server:

```bash
./build/Myapp/Myapp.exe
```

---

## CLI Usage

### View Available Commands

```bash
./goduck.exe --help
```

### Build a Project

Recursively parses, type-checks, and compiles all `.duck` files in the specified directory:

```bash
./goduck.exe build <directory-or-file>
```

### Scaffolding

Scaffold a complete NestJS-like starter project structure:

```bash
./goduck.exe new <project-name>
```

```
# 1. Delete the installed goduck executable
rm -f ~/go/bin/goduck.exe

# 2. Delete the .goduck configuration folder (resets first-time setup markers)
rm -rf ~/.goduck

# 3. Clear Go's module cache to force it to download the updated tag from GitHub
go clean -modcache
```

```
GOPROXY=direct go install github.com/shishir1290/goduck/cmd/goduck@v0.0.1

```
