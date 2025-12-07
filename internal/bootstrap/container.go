package bootstrap

import (
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/pkg/database"
)

// Container holds all dependencies
type Container struct {
	// Repositories
	UserRepo repositories.UserRepository

	// Services
	AuthService services.AuthService

	// Controllers
	AuthController *controllers.AuthController
}

// NewContainer initializes all dependencies
func NewContainer() *Container {
	db := database.GetDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)

	return &Container{
		UserRepo:       userRepo,
		AuthService:    authService,
		AuthController: authController,
	}
}

// Global container instance
var container *Container

// GetContainer returns the global container
func GetContainer() *Container {
	if container == nil {
		container = NewContainer()
	}
	return container
}

// InitContainer initializes the container (called after database is ready)
func InitContainer() {
	container = NewContainer()
}

// GetAuthController returns the auth controller
func (c *Container) GetAuthController() *controllers.AuthController {
	return c.AuthController
}
