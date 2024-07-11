package router

import (
	_ "PoliticianRating/docs"
	"PoliticianRating/pkg/controllers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(r *gin.Engine) {
	ratingRoutes := r.Group("/rating")
	{
		ratingRoutes.GET("/", controllers.GetAllRatings)
		ratingRoutes.GET("/:id", controllers.GetRating)
		ratingRoutes.PUT("/:id/:action", controllers.UpdateRating)
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

}
