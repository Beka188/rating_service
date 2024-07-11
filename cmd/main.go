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

	r := gin.Default()
	r.GET("/rating", func(c *gin.Context) {
		sortBy := c.DefaultQuery("sortBy", "+date")
		allRatings, err := ratingRepo.GetAllRatings(sortBy)
		if err != nil {
			if err.Error() == fmt.Sprintf("invalid order: %s", sortBy) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			c.JSON(http.StatusInternalServerError, gin.H{})
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

	r.PUT("/rating/:id/:action", func(c *gin.Context) {
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
		case "decrement":
			newRating, err = ratingRepo.UpdateUserRating(id, false)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
			return
		}
		if err != nil {
			if err.Error() == fmt.Sprintf("user with ID %d does not exist", id) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, newRating)
	})
	r.Run() // listen and serve on 0.0.0.0:8080

}
