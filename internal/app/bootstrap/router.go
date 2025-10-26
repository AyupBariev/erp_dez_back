package bootstrap

import (
	"erp/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h Handlers) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(middleware.RecoveryWithLog())
	// Public routes

	// Protected routes
	protected := router.Group("/api")
	protected.POST("/login", h.Auth.LoginHandler)
	protected.Use(middleware.AuthMiddleware(h.Auth.AuthService))
	{
		users := protected.Group("/users")
		{
			users.POST("/", h.User.CreateUser) // POST /api/users
			//users.GET("/", h.User.ListUsers)        // GET /api/users
			//users.GET("/:id", h.User.GetUser)       // GET /api/users/:id
			//users.PUT("/:id", h.User.UpdateUser)    // PUT /api/users/:id
			//users.DELETE("/:id", h.User.DeleteUser) // DELETE /api/users/:id
		}

		// Orders group
		orders := protected.Group("/orders")
		{
			orders.POST("/", h.Order.CreateOrderHandler)             // POST /api/orders
			orders.POST("/assign-order", h.Order.AssignOrderHandler) // POST /api/orders
			orders.GET("/", h.Order.ListOrders)                      // GET /api/orders
			//	orders.GET("/:id", h.Order.GetOrder)       // GET /api/orders/:id
			//	orders.PUT("/:id", h.Order.UpdateOrder)    // PUT /api/orders/:id
			//	orders.DELETE("/:id", h.Order.DeleteOrder) // DELETE /api/orders/:id
		}
		//
		engineers := protected.Group("/engineers")
		{
			engineers.POST("/", h.Engineer.CreateEngineer) // POST /api/engineers
			engineers.GET("/", h.Engineer.ListEngineers)   // GET /api/engineers
			//engineers.GET("/:id", h.Engineer.GetEngineer)       // GET /api/engineers/:id
			//engineers.PUT("/:id", h.Engineer.UpdateEngineer)    // PUT /api/engineers/:id
			//engineers.DELETE("/:id", h.Engineer.DeleteEngineer) // DELETE /api/engineers/:id

			engineers.POST("/accept-engineer", h.Admin.ApproveEngineer)
		}

		// Допустим выдача заказа
		//orders.POST("/:id/assign/:engineerId", h.Order.AssignOrderToEngineer)

		protected.POST("/logout", h.Auth.LogoutHandler)
	}

	return router
}
