package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new Goduck project structure",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		projectDir := projectName

		fmt.Printf("Creating a new Goduck project: %s...\n", projectName)

		// Create directories
		dirs := []string{
			filepath.Join(projectDir, "src"),
			filepath.Join(projectDir, "src", "users"),
			filepath.Join(projectDir, "src", "auth"),
		}

		for _, d := range dirs {
			if err := os.MkdirAll(d, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", d, err)
			}
		}

		// Capitalize app name
		appName := strings.Title(filepath.Base(projectDir))

		// Create files
		files := map[string]string{
			filepath.Join(projectDir, "src", "main.duck"): fmt.Sprintf(`import { AppModule } from "./app.module.duck";
import { AuthModule } from "./auth.module.duck";
import { UserModule } from "./users/user.module.duck";

app %s {
    port: number = 9000
    apiPath: string = "/api"
}
`, appName),

			filepath.Join(projectDir, "src", "app.module.duck"): `import { AppController } from "./app.controller.duck";
import { AppService } from "./app.service.duck";

@Module
func AppModule {
    controller: AppController;
    service: AppService;
}
`,

			filepath.Join(projectDir, "src", "auth", "auth.module.duck"): `import { AuthController } from "./auth.controller.duck";
import { AuthService } from "./auth.service.duck";

@Module
func AuthModule {
    controller: AuthController;
    service: AuthService;
}
`,

			filepath.Join(projectDir, "src", "users", "user.module.duck"): `import { UsersController } from "./users.controller.duck";
import { UsersService } from "./users.service.duck";
import { UsersRepository } from "./users.repository.duck";

@Module
func UserModule {
    controller: UsersController;
    service: UsersService;
    repo: UsersRepository;
}
`,

			filepath.Join(projectDir, "src", "app.controller.duck"): `import { AppService } from "./app.service.duck";

@Controller("/")
func AppController {
    service: AppService;

    @GET("/health")
    async health(): string {
        return this.service.getHealth();
    }
}
`,

			filepath.Join(projectDir, "src", "app.service.duck"): `@Service("app")
service AppService {
    async getHealth(): string {
        return "OK";
    }
}
`,

			filepath.Join(projectDir, "src", "users", "user.duck"): `func User {
    id: number;
    name: string;
    email: string;
    active: boolean;
}
`,

			filepath.Join(projectDir, "src", "users", "users.dto.duck"): `func CreateUserDto {
    name: string;
    email: string;
}
`,

			filepath.Join(projectDir, "src", "users", "users.repository.duck"): `@Service("usersRepository")
func UsersRepository {
    async save(user: User): string {
        return "saved";
    }
}
`,

			filepath.Join(projectDir, "src", "users", "users.service.duck"): `import { UsersRepository } from "./users.repository.duck";
import { CreateUserDto } from "./users.dto.duck";
import { User } from "./user.duck";

@Service("user")
service UsersService {
    repo: UsersRepository;

    async findAll(): string {
        return "all users";
    }

    async findOne(id: string): string {
        return "user with id " + id;
    }

    async create(dto: CreateUserDto): string {
        return "created " + dto.name;
    }
}
`,

			filepath.Join(projectDir, "src", "users", "users.controller.duck"): `import { UsersService } from "./users.service.duck";

@Controller("/users")
func UsersController {
    service: UsersService;

    @GET("/")
    async index(): string {
        return this.service.findAll();
    }

    @GET("/:id")
    async show(id: string): string {
        return this.service.findOne(id);
    }
}
`,

			filepath.Join(projectDir, "src", "auth", "auth.guard.duck"): `@Service("authGuard")
func AuthGuard {
    async canActivate(): boolean {
        return true;
    }
}
`,

			filepath.Join(projectDir, "src", "auth", "auth.service.duck"): `@Service("auth")
service AuthService {
    async validateUser(username: string): boolean {
        return true;
    }
}
`,

			filepath.Join(projectDir, "src", "auth", "auth.controller.duck"): `import { AuthService } from "./auth.service.duck";

@Controller("/auth")
func AuthController {
    authService: AuthService;

    @POST("/login")
    async login(): string {
        return "logged in";
    }
}
`,

			filepath.Join(projectDir, ".gitignore"): `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
build/

# Dependency directories
vendor/

# Environment files
.env
.env.local
.env.*.local

# IDEs
.idea/
.vscode/
*.suo
*.ntvs*
*.njsproj
*.sln
*.sw?
`,

			filepath.Join(projectDir, ".env.example"): `PORT=8080
DB_HOST=localhost
DB_PORT=27017
DB_NAME=goduck
`,

			filepath.Join(projectDir, "README.md"): fmt.Sprintf(`# %s: Goduck Backend Project

This project is built using the Goduck Backend DSL framework, generating highly optimized, dependency-injection-enabled Go web applications.

## Directory Structure

§§§
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
§§§

## CLI Commands

### 1. Build the Project
Compile the §.duck§ source code and generate the final Go application & executable:
§§§bash
goduck build src
§§§

### 2. Run the Server

**Development (with Hot-Reloading):**
§§§bash
goduck run
§§§

**Production (run the compiled binary directly):**
§§§bash
./build/%s/%s.exe
§§§

### 3. Create a New Module
Generate a new controller, service, or model module skeleton:
§§§bash
goduck g module <module-name>
§§§

### 4. Database Migrations
Manage database schemas and run Prisma migrations:
§§§bash
goduck db migrate dev --name init
goduck db push
§§§
`, appName, appName, appName),
		}

		for path, content := range files {
			content = strings.ReplaceAll(content, "§§§", "```")
			content = strings.ReplaceAll(content, "§", "`")
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", path, err)
			}
		}

		fmt.Println("✓ Project created successfully!")
		fmt.Println("\nTo build and run the project:")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Println("  goduck build src")
		return nil
	},
}
