//4.4. Score Service
//Функциональность: Расчёт и обновление кармы.
//API: /ratings, /ratings/{politician_id}

package main

import (
	"PoliticianRating/pkg/model"
	"PoliticianRating/pkg/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	_ "net/http"
	"strconv"
)

func main() {
	ratingRepo := repository.NewRatingRepository()

	//dropTableSQL := `DROP TABLE IF EXISTS Rating;`
	//
	//_, err := ratingRepo.Db.Exec(dropTableSQL)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//rat, err := ratingRepo.GetUserRating(4)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(rat)
	//ratingRepo.InsertRating(3, 10)
	//ratingRepo.DecrementRating(2)
	//fmt.Println(time.Now())
	//fmt.Println(ratingRepo.GetAllRatings(""))
	//
	r := gin.Default()
	r.GET("/rating", func(c *gin.Context) {
		allRatings, err := ratingRepo.GetAllRatings("asc")
		if err != nil {
			fmt.Println(err)
			return
		}
		c.JSON(http.StatusOK, allRatings)
	})

	r.GET("/rating/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be an integer"})
			return
		}
		rating, err := ratingRepo.GetUserRating(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rating)

	})

	r.GET("/rating/:id/:action", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be an integer"})
		}
		action := c.Param("action")
		var newRating model.Rating
		switch action {
		case "increment":
			newRating, err = ratingRepo.UpdateUserRating(id, true)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		case "decrement":
			newRating, err = ratingRepo.UpdateUserRating(id, false)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		}
		c.JSON(http.StatusOK, newRating)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	//for _, row := range allRatings {
	//	fmt.Println(row.ID, row.UserID, row.Score)
	//}

}
