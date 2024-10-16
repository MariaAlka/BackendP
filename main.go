package main

import (
	"errors"
	"fmt"
	"intership/controllers"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-michi/michi"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// Middleware to handle CORS
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Allow your frontend
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials (cookies, tokens)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Print the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", cwd)

	// Path to .env file in the root directory
	envPath := filepath.Join(cwd, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Fatalf(".env file does not exist at path: %s", envPath)
	}

	// Load .env file
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Check if environment variables are loaded correctly
	fmt.Println("MIGRATIONS_ROOT:", os.Getenv("MIGRATIONS_ROOT"))
	fmt.Println("DATABASE_URL:", os.Getenv("DATABASE_URL"))

	// Connect to the database
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Perform database migrations
	migrationsRoot := os.Getenv("MIGRATIONS_ROOT")
	absMigrationsPath := GetRootPath(migrationsRoot)
	migrationPath := "file://" + strings.ReplaceAll(filepath.ToSlash(absMigrationsPath), "\\", "/")
	fmt.Println("Migrations path:", migrationPath)

	mig, err := migrate.New(
		migrationPath,
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := mig.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	} else if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrations: %s", err.Error())
	}

	// Set the database connection in the controllers
	controllers.SetDB(db)

	// Setup router and routes
	r := michi.NewRouter()

	// Serve uploaded files
	r.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// User routes
	r.Route("/", func(sub *michi.Router) {
		// User CRUD routes
		sub.HandleFunc("GET users", controllers.IndexUserHandler)          // GET /users
		sub.HandleFunc("GET users/{id}", controllers.ShowUserHandler)      // GET /users/{id}
		sub.HandleFunc("PUT users/{id}", controllers.UpdateUserHandler)    // PUT /users/{id}
		sub.HandleFunc("DELETE users/{id}", controllers.DeleteUserHandler) // DELETE /users/{id}
		sub.HandleFunc("POST users/signup", controllers.SignUpHandler)     // POST /users/signup
		sub.HandleFunc("POST users/login", controllers.LoginHandler)       // POST /users/login

		// Vendor admin routes
		sub.HandleFunc("POST vendor_admins", controllers.CreateVendorAdminHandler)                         // POST /vendor_admins
		sub.HandleFunc("GET vendor_admins", controllers.IndexVendorAdminsHandler)                          // GET /vendor_admins
		sub.HandleFunc("GET vendor_admins/{user_id}/{vendor_id}", controllers.ShowVendorAdminHandler)      // GET /vendor_admins/{user_id}/{vendor_id}
		sub.HandleFunc("DELETE vendor_admins/{user_id}/{vendor_id}", controllers.DeleteVendorAdminHandler) // DELETE /vendor_admins/{user_id}/{vendor_id}
		sub.HandleFunc("PUT vendor_admins/{user_id}/{vendor_id}", controllers.UpdateVendorAdminHandler)

		// Vendor routes
		sub.HandleFunc("GET vendors", controllers.IndexVendorHandler)          // GET /vendors
		sub.HandleFunc("GET vendors/{id}", controllers.ShowVendorHandler)      // GET /vendors/{id}
		sub.HandleFunc("PUT vendors/{id}", controllers.UpdateVendorHandler)    // PUT /vendors/{id}
		sub.HandleFunc("DELETE vendors/{id}", controllers.DeleteVendorHandler) // DELETE /vendors/{id}
		sub.HandleFunc("POST vendors/signup", controllers.SignUpVendorHandler) // POST /vendors/signup

		// User roles routes
		sub.HandleFunc("GET user_roles", controllers.IndexUserRolesHandler)                        // GET /user_roles
		sub.HandleFunc("GET user_roles/{user_id}/{role_id}", controllers.ShowUserRoleHandler)      // GET /user_roles/{user_id}/{role_id}
		sub.HandleFunc("POST user_roles", controllers.CreateUserRoleHandler)                       // POST /user_roles
		sub.HandleFunc("DELETE user_roles/{user_id}/{role_id}", controllers.DeleteUserRoleHandler) // DELETE /user_roles/{user_id}/{role_id}
		sub.HandleFunc("PUT user_roles", controllers.UpdateUserRoleHandler)                        // PUT /user_roles/{user_id}

		// Item routes
		sub.HandleFunc("POST items", controllers.CreateItemHandler)        // POST /items
		sub.HandleFunc("GET items", controllers.IndexItemHandler)          // GET /items
		sub.HandleFunc("GET items/{id}", controllers.ShowItemHandler)      // GET /items/{id}
		sub.HandleFunc("PUT items/{id}", controllers.UpdateItemHandler)    // PUT /items/{id}
		sub.HandleFunc("DELETE items/{id}", controllers.DeleteItemHandler) // DELETE /items/{id}
		//tables routes
		sub.HandleFunc("GET tables", controllers.IndexTableHandler)          // GET /tables
		sub.HandleFunc("GET tables/{id}", controllers.ShowTableHandler)      // GET /tables/{id}
		sub.HandleFunc("POST tables", controllers.CreateTableHandler)        // POST /tables
		sub.HandleFunc("PUT tables/{id}", controllers.UpdateTableHandler)    // PUT /tables/{id}
		sub.HandleFunc("DELETE tables/{id}", controllers.DeleteTableHandler) // DELETE /tables/{id}

		//ordersrouts
		sub.HandleFunc("POST orders", controllers.CreateOrderHandler)        // POST /orders
		sub.HandleFunc("GET orders", controllers.IndexOrderHandler)          // GET /orders
		sub.HandleFunc("GET orders/{id}", controllers.ShowOrderHandler)      // GET /orders/{id}
		sub.HandleFunc("PUT orders/{id}", controllers.UpdateOrderHandler)    // PUT /orders/{id}
		sub.HandleFunc("DELETE orders/{id}", controllers.DeleteOrderHandler) // DELETE /orders/

		// Order Items CRUD routes
		sub.HandleFunc("POST order_items", controllers.CreateOrderItemHandler)        // POST /order_items
		sub.HandleFunc("GET order_items", controllers.IndexOrderItemHandler)          // GET /order_items
		sub.HandleFunc("GET order_items/{id}", controllers.ShowOrderItemHandler)      // GET /order_items/{id}
		sub.HandleFunc("PUT order_items/{id}", controllers.UpdateOrderItemHandler)    // PUT /order_items/{id}
		sub.HandleFunc("DELETE order_items/{id}", controllers.DeleteOrderItemHandler) // DELETE /order_items/{id}

		// Carts CRUD routes
		sub.HandleFunc("POST carts", controllers.CreateCartHandler)        // POST /carts
		sub.HandleFunc("GET carts", controllers.IndexCartHandler)          // GET /carts
		sub.HandleFunc("GET carts/{id}", controllers.ShowCartHandler)      // GET /carts/{id}
		sub.HandleFunc("PUT carts/{id}", controllers.UpdateCartHandler)    // PUT /carts/{id}
		sub.HandleFunc("DELETE carts/{id}", controllers.DeleteCartHandler) // DELETE /carts/{id}

		sub.HandleFunc("GET cart_items", controllers.IndexCartItemsHandler)                        // GET /cart_items
		sub.HandleFunc("GET cart_items/{cart_id}/{item_id}", controllers.ShowCartItemHandler)      // GET /cart_items/{cart_id}/{item_id}
		sub.HandleFunc("POST cart_items", controllers.CreateCartItemHandler)                       // POST /cart_items
		sub.HandleFunc("PUT cart_items/{cart_id}/{item_id}", controllers.UpdateCartItemHandler)    // PUT /cart_items/{cart_id}/{item_id}
		sub.HandleFunc("DELETE cart_items/{cart_id}/{item_id}", controllers.DeleteCartItemHandler) // DELETE /cart_items/{cart_id}/{item_id}

	})

	// Wrap the router with the CORS middleware
	corsRouter := enableCors(r)

	// Start the server with CORS-enabled routes
	fmt.Println("Starting server on port 8000")
	if err := http.ListenAndServe(":8000", corsRouter); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// GetRootPath resolves the absolute path of a given directory relative to the project root
func GetRootPath(dir string) string {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}
	return absPath
}
