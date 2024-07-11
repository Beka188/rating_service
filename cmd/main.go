//4.4. Score Service
//Функциональность: Расчёт и обновление кармы.
//API: /ratings, /ratings/{politician_id}

package main

import (
	_ "PoliticianRating/docs"
	"PoliticianRating/internal/router"
	"PoliticianRating/pkg/database"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// @title           Rating Service API
// @version         1.0
// @description		This server is part of the Politician application, which allows users to view and update ratings for politicians. It uses the Gin framework for handling HTTP requests and provides endpoints for retrieving all ratings, retrieving a specific rating by ID, and updating a rating based on an action (increment or decrement)
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

func main() {
	errDb := database.InitDB()
	if errDb != nil {
		panic(errDb)
	}
	r := gin.Default()
	router.InitRouter(r)
	r.Run(":8080")
}
