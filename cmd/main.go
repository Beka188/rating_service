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

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	errDb := database.InitDB()
	if errDb != nil {
		panic(errDb)
	}
	r := gin.Default()
	router.InitRouter(r)
	r.Run(":8080")
}
