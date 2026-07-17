package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/api/handlers"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/api/middleware"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/repository"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/service"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/config"
	"github.com/Raunakpratapkushwaha/Batwara/backend/pkg/token"
)

// SetupRouter initializes all repositories, services, and handlers,
// registers public/protected routes, and returns the configured Gin engine.
func SetupRouter(db *sqlx.DB, tokenMaker *token.TokenMaker, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// -------------------------------------------------------------------------
	// 1. Initialize Repositories (Database Layer)
	// -------------------------------------------------------------------------
	// Passing db.DB (which is *sql.DB) to repositories expecting standard SQL
	userRepo := repository.NewSQLUserRepository(db.DB)
	groupRepo := repository.NewGroupRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// -------------------------------------------------------------------------
	// 2. Initialize Services (Business Logic Layer)
	// -------------------------------------------------------------------------
	authService := service.NewAuthService(userRepo, tokenMaker, cfg)
	groupService := service.NewGroupService(groupRepo)
	expenseService := service.NewExpenseService(expenseRepo, groupRepo)

	// Inject the expenseService to let the algorithm fetch updated net balances
	settlementService := service.NewSettlementService(expenseService)

	// Inject both payment and group repositories for permission checking
	paymentService := service.NewPaymentService(paymentRepo, groupRepo)

	// -------------------------------------------------------------------------
	// 3. Initialize Handlers (Controller Layer)
	// -------------------------------------------------------------------------
	authHandler := handlers.NewAuthHandler(authService)
	groupHandler := handlers.NewGroupHandler(groupService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	settlementHandler := handlers.NewSettlementHandler(settlementService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// -------------------------------------------------------------------------
	// 4. Register Routes
	// -------------------------------------------------------------------------

	// Public Authentication Routes
	r.POST("/api/v1/auth/register", authHandler.Register)
	r.POST("/api/v1/auth/login", authHandler.Login)

	// Protected Routes (Guarded by JWT Auth Middleware)
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(tokenMaker))
	{
		// --- Group Management ---
		protected.POST("/groups", groupHandler.CreateGroup)
		protected.GET("/groups", groupHandler.ListGroups)
		protected.POST("/groups/:id/members", groupHandler.AddMember)

		// --- Core Expense & Splits ---
		protected.POST("/expenses", expenseHandler.CreateExpense)
		protected.GET("/groups/:id/expenses", expenseHandler.GetGroupExpenses)
		protected.GET("/groups/:id/balances", expenseHandler.GetGroupBalances)

		// --- Debt Simplification ---
		protected.GET("/groups/:id/settlements", settlementHandler.GetGroupSettlements)

		// --- Settlement Payment Tracking ---
		protected.POST("/groups/:id/payments", paymentHandler.RecordPayment)
		protected.GET("/groups/:id/payments", paymentHandler.GetGroupPayments)
	}

	return r
}
