package controllers

import (
	"PoliticianRating/pkg/model"
	"PoliticianRating/pkg/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetAllRatings(c *gin.Context) {
	sortBy := c.DefaultQuery("sortBy", "date")
	allRatings, err := services.GetAllRatings(sortBy)
	if err != nil {
		if err.Error() == fmt.Sprintf("invalid order: %s", sortBy) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, allRatings)
}

func GetRating(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be an integer"})
		return
	}
	rating, err := services.GetUserRating(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rating)

}

func UpdateRating(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be an integer"})
	}
	action := c.Param("action")
	var newRating model.Rating
	switch action {
	case "increment":
		newRating, err = services.UpdateUserRating(id, true)
	case "decrement":
		newRating, err = services.UpdateUserRating(id, false)
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
}
