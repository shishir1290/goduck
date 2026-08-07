# Test: Goduck Backend Project

This project is built using the Goduck Backend DSL framework, generating highly optimized, dependency-injection-enabled Go web applications.

## Directory Structure

```
src/
├── app.controller.duck
├── app.duck
├── app.service.duck
├── main.duck
├── auth/
│   ├── auth.controller.duck
│   ├── auth.service.duck
│   └── auth.guard.duck
└── users/
    ├── user.duck
    ├── users.dto.duck
    ├── users.repository.duck
    ├── users.service.duck
    └── users.controller.duck
```

## CLI Commands

### 1. Build the Project
Compile the `.duck` source code and generate the final Go application & executable:
```bash
goduck build src
```

### 2. Run the Server
Build and run the compiled binary directly:
```bash
./build/Test/Test.exe
```

### 3. Create a New Module
Generate a new controller, service, or model module skeleton:
```bash
goduck g module <module-name>
```

### 4. Database Migrations
Manage database schemas and run Prisma migrations:
```bash
goduck db migrate dev --name init
goduck db push
```
