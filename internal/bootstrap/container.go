package bootstrap

import (
	"context"
	"log"

	"github.com/habbazettt/nutrisnap-server/config"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/internal/workers"
	"github.com/habbazettt/nutrisnap-server/pkg/database"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
	"github.com/habbazettt/nutrisnap-server/pkg/oauth"
	"github.com/habbazettt/nutrisnap-server/pkg/openfoodfacts"
	"github.com/habbazettt/nutrisnap-server/pkg/storage"
)

// Container holds all dependencies
type Container struct {
	// JWT
	JWTManager *jwt.Manager

	// OAuth
	GoogleOAuth *oauth.GoogleOAuth

	// Storage
	StorageClient *storage.Client

	// External APIs
	OFFClient *openfoodfacts.Client

	// Repositories
	UserRepo       repositories.UserRepository
	ScanRepo       repositories.ScanRepository
	ProductRepo    repositories.ProductRepository
	CorrectionRepo repositories.CorrectionRepository

	// Services
	AuthService    services.AuthService
	UserService    services.UserService
	AdminService   services.AdminService
	ScanService    services.ScanService
	ProductService services.ProductService
	OCRService     services.OCRService

	// Workers
	OCRWorker *workers.OCRWorker

	// Controllers
	AuthController       *controllers.AuthController
	UserController       *controllers.UserController
	AdminController      *controllers.AdminController
	ScanController       *controllers.ScanController
	ProductController    *controllers.ProductController
	CorrectionController *controllers.CorrectionController
	CompareController    *controllers.CompareController
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

	// Initialize MinIO storage client
	storageClient, err := storage.NewClient(storage.Config{
		Endpoint:  cfg.MinIO.Endpoint,
		PublicURL: cfg.MinIO.PublicURL,
		AccessKey: cfg.MinIO.AccessKey,
		SecretKey: cfg.MinIO.SecretKey,
		Bucket:    cfg.MinIO.Bucket,
		UseSSL:    cfg.MinIO.UseSSL,
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize storage client: %v", err)
	} else {
		// Ensure bucket exists
		if err := storageClient.EnsureBucket(context.Background()); err != nil {
			log.Printf("Warning: Failed to ensure bucket: %v", err)
		}
	}

	// Initialize OpenFoodFacts client
	offClient := openfoodfacts.NewClient()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	scanRepo := repositories.NewScanRepository(db)
	productRepo := repositories.NewProductRepository(db)
	correctionRepo := repositories.NewCorrectionRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtManager, googleOAuth)
	userService := services.NewUserService(userRepo)
	adminService := services.NewAdminService(userRepo)
	productService := services.NewProductService(productRepo, offClient)
	ocrService := services.NewOCRService(storageClient)

	// Initialize Workers
	ocrWorker := workers.NewOCRWorker(scanRepo, productRepo, ocrService, 100) // Buffer 100 jobs

	// ScanService needs ScanQueue (implemented by ocrWorker)
	scanService := services.NewScanService(scanRepo, storageClient, productService, ocrWorker)

	// Initialize Correction Service
	correctionService := services.NewCorrectionService(correctionRepo, scanRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	adminController := controllers.NewAdminController(adminService)
	scanController := controllers.NewScanController(scanService)
	productController := controllers.NewProductController(productService)
	correctionController := controllers.NewCorrectionController(correctionService)

	// Initialize Compare Service and Controller
	compareService := services.NewCompareService(productRepo, scanRepo)
	compareController := controllers.NewCompareController(compareService)

	return &Container{
		JWTManager:           jwtManager,
		GoogleOAuth:          googleOAuth,
		StorageClient:        storageClient,
		OFFClient:            offClient,
		UserRepo:             userRepo,
		ScanRepo:             scanRepo,
		ProductRepo:          productRepo,
		CorrectionRepo:       correctionRepo,
		AuthService:          authService,
		UserService:          userService,
		AdminService:         adminService,
		ScanService:          scanService,
		ProductService:       productService,
		OCRService:           ocrService,
		OCRWorker:            ocrWorker,
		AuthController:       authController,
		UserController:       userController,
		AdminController:      adminController,
		ScanController:       scanController,
		ProductController:    productController,
		CorrectionController: correctionController,
		CompareController:    compareController,
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

// GetAdminController returns the admin controller
func (c *Container) GetAdminController() *controllers.AdminController {
	return c.AdminController
}

// GetScanController returns the scan controller
func (c *Container) GetScanController() *controllers.ScanController {
	return c.ScanController
}

// GetProductController returns the product controller
func (c *Container) GetProductController() *controllers.ProductController {
	return c.ProductController
}

// GetCorrectionController returns the correction controller
func (c *Container) GetCorrectionController() *controllers.CorrectionController {
	return c.CorrectionController
}

// GetCompareController returns the compare controller
func (c *Container) GetCompareController() *controllers.CompareController {
	return c.CompareController
}

// GetJWTManager returns the JWT manager
func (c *Container) GetJWTManager() *jwt.Manager {
	return c.JWTManager
}
