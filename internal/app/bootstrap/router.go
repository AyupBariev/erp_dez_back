package bootstrap

import (
	"erp/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h Handlers) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())           // базовый recover от Gin
	router.Use(middleware.ErrorLogger()) // наш кастомный
	// Public routes

	// Protected routes
	prefixApi := router.Group("/api")
	prefixApi.POST("/login", h.Auth.LoginHandler)

	prefixApi.GET("/reports/link/:token", h.Report.GetReportByToken)
	prefixApi.POST("/reports/submit", h.Report.SubmitReport)
	prefixApi.Use(middleware.AuthMiddleware(h.Auth.AuthService))
	{
		users := prefixApi.Group("/users")
		{
			users.POST("/", h.User.CreateUser) // POST /api/users
			//users.GET("/", h.User.ListUsers)        // GET /api/users
			//users.GET("/:id", h.User.GetUser)       // GET /api/users/:id
			//users.PUT("/:id", h.User.UpdateUser)    // PUT /api/users/:id
			//users.DELETE("/:id", h.User.DeleteUser) // DELETE /api/users/:id
		}

		// ORDERS
		orders := prefixApi.Group("/orders")
		{
			orders.POST("/", h.Order.CreateOrderHandler) // POST /api/orders
			orders.GET("/", h.Order.ListOrders)          // GET /api/orders
			//	orders.GET("/:id", h.Order.GetOrder)       // GET /api/orders/:id
			orders.PUT("/:erp_number", h.Order.UpdateOrder)    // PUT /api/orders/:id
			orders.DELETE("/:erp_number", h.Order.DeleteOrder) // DELETE /api/orders/:id

			orders.POST("/assign", h.Order.AssignOrderHandler)
			orders.POST("/unassign", h.Order.UnAssignOrderHandler)
		}

		// REPORTS
		reports := prefixApi.Group("/reports")
		{
			reports.GET("/cash", h.Report.ListCashReports) // ?from=..&to=..
			reports.POST("/cash/:erp_number/receive", h.Report.ReceiveCash)
		}

		//REPEAT REQUEST
		repeat := prefixApi.Group("/repeat-requests")
		{
			repeat.GET("/", h.RepeatRequest.List)
			repeat.POST("/:id/confirm", h.RepeatRequest.Confirm)
		}

		agg := prefixApi.Group("/aggregator")
		{
			agg.GET("/daily", h.AggregatorPayout.ListDayPayouts)
		}

		profit := prefixApi.Group("/profit")
		{
			profit.GET("/", h.Profit.GetProfit)
		}

		engineers := prefixApi.Group("/engineers")
		{
			engineers.POST("/", h.Engineer.CreateEngineer) // POST /api/engineers
			engineers.GET("/", h.Engineer.ListEngineers)   // GET /api/engineers
			//engineers.GET("/:id", h.Engineer.GetEngineer)       // GET /api/engineers/:id
			//engineers.PUT("/:id", h.Engineer.UpdateEngineer)    // PUT /api/engineers/:id
			//engineers.DELETE("/:id", h.Engineer.DeleteEngineer) // DELETE /api/engineers/:id
			engineers.POST("/accept-engineer", h.Admin.ApproveEngineer)
			engineers.POST("/:engineer_id/working-status", h.Admin.ToggleWorkingStatusEngineer)

			engineers.GET("/month-payouts", h.Engineer.GetMonthPayouts)
			engineers.POST("/advance", h.Engineer.PayAdvance)
		}

		motivations := prefixApi.Group("/motivations")
		{
			motivations.GET("/engineer", h.EngineerMotivation.GetMonthlyMotivation)

		}

		dictionaries := prefixApi.Group("/dictionaries")
		{
			dictionaries.GET("/aggregators", h.DictHandler.HandleDictionary("aggregators"))
			dictionaries.GET("/aggregators/:id", h.DictHandler.HandleDictionary("aggregators"))
			dictionaries.POST("/aggregators", h.DictHandler.HandleDictionary("aggregators"))
			dictionaries.PUT("/aggregators/:id", h.DictHandler.HandleDictionary("aggregators"))
			dictionaries.DELETE("/aggregators/:id", h.DictHandler.HandleDictionary("aggregators"))

			dictionaries.GET("/problems", h.DictHandler.HandleDictionary("problems"))
			dictionaries.GET("/problems/:id", h.DictHandler.HandleDictionary("problems"))
			dictionaries.POST("/problems", h.DictHandler.HandleDictionary("problems"))
			dictionaries.PUT("/problems/:id", h.DictHandler.HandleDictionary("problems"))
			dictionaries.DELETE("/problems/:id", h.DictHandler.HandleDictionary("problems"))
		}

		prefixApi.POST("/logout", h.Auth.LogoutHandler)
	}

	return router
}
