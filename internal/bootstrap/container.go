package bootstrap

import (
	"github.com/habbazettt/nutrisnap-server/config"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/pkg/database"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
	"github.com/habbazettt/nutrisnap-server/pkg/oauth"
)

// Container holds all dependencies
type Container struct {
	// JWT
	JWTManager *jwt.Manager

	// OAuth
	GoogleOAuth *oauth.GoogleOAuth

	// Repositories
	UserRepo repositories.UserRepository

	// Services
	AuthService services.AuthService
	UserService services.UserService

	// Controllers
	AuthController *controllers.AuthController
	UserController *controllers.UserController
}

// NewContainer initializes all dependencies
func NewContainer() *Container {
	db := database.GetDB()
	cfg := config.Get()

	// Initialize JWT manager
	jwtManager := jwt.NewManager(jwt.Config{
		SecretKey:     cfg.JWT.Secret,
		AccessExpiry:  cfg.JWT.AccessExpiry,
		RefreshExpiry: cfg.JWT.RefreshExpiry,
		Issuer:        cfg.JWT.Issuer,
	})

	// Initialize Google OAuth
	googleOAuth := oauth.NewGoogleOAuth(oauth.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		RedirectURL:  cfg.Google.RedirectURL,
	})

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtManager, googleOAuth)
	userService := services.NewUserService(userRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	return &Container{
		JWTManager:     jwtManager,
		GoogleOAuth:    googleOAuth,
		UserRepo:       userRepo,
		AuthService:    authService,
		UserService:    userService,
		AuthController: authController,
		UserController: userController,
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

// GetUserController returns the user controller
func (c *Container) GetUserController() *controllers.UserController {
	return c.UserController
}

// GetJWTManager returns the JWT manager
func (c *Container) GetJWTManager() *jwt.Manager {
	return c.JWTManager
}
