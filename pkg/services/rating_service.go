package services

import (
	"PoliticianRating/pkg/database"
	"PoliticianRating/pkg/model"
	"fmt"
)

// GetAllRatings godoc
// @Summary      Show all ratings
// @Description  get all ratings, optionally can sort by date (date or -date), and score (score or +score)
// @Accept       json
// @Produce      json
// @Param sortBy query string false "Sort by"
// @Success      200  {object}  model.Rating
// @Router       /rating/ [get]
func GetAllRatings(order string) ([]model.Rating, error) {
	var ratings []model.Rating
	var orderBy string
	fmt.Println(order)
	switch order {
	case "score":
		orderBy = "ORDER BY Score ASC"
	case "-score":
		orderBy = "ORDER BY Score DESC"
	case "date":
		orderBy = "ORDER BY UpdatedAt ASC"
	case "-date":
		orderBy = "ORDER BY UpdatedAt DESC"
	default:
		return nil, fmt.Errorf("invalid order: %s", order)
	}

	selectAllSQL := fmt.Sprintf("SELECT ID, UserId, Score, CreatedAt, UpdatedAt FROM rating %s", orderBy)
	rows, err := database.DB.Query(selectAllSQL)
	if err != nil {
		return nil, fmt.Errorf("error querying ratings: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rating model.Rating
		err := rows.Scan(&rating.ID, &rating.UserID, &rating.Score, &rating.CreatedAt, &rating.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning rating row: %v", err)
		}
		ratings = append(ratings, rating)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over ratings: %v", err)
	}
	return ratings, nil
}

func GetUserRating(userID int) (model.Rating, error) {
	var rating model.Rating
	selectUserSql := fmt.Sprintf("SELECT * FROM rating WHERE UserID = %d", userID)
	row, err := database.DB.Query(selectUserSql)
	if err != nil {
		return rating, fmt.Errorf("error querying ratings: %v", err)
	}
	if !row.Next() {
		return rating, fmt.Errorf("no rating found for user ID %d", userID)
	}
	err2 := row.Scan(&rating.ID, &rating.UserID, &rating.Score, &rating.CreatedAt, &rating.UpdatedAt)
	if err2 != nil {
		return rating, fmt.Errorf("error scanning ratings: %v", err2)
	}
	return rating, nil
}

// UpdateUserRating godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  model.Rating
// @Router       /rating/{id} [put]
func UpdateUserRating(userID int, increment bool) (model.Rating, error) {
	var newRating model.Rating

	var userExists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM rating WHERE UserID = ?)", userID).Scan(&userExists)
	if err != nil {
		return newRating, fmt.Errorf("error checking if user exists: %v", err)
	}
	if !userExists {
		return newRating, fmt.Errorf("user with ID %d does not exist", userID)
	}

	var updateScoreSQL string
	if increment {
		updateScoreSQL = "UPDATE rating SET Score = Score + 1 WHERE UserID = ?"
	} else {
		updateScoreSQL = "UPDATE rating SET Score = Score - 1 WHERE UserID = ?"
	}

	stmt, err := database.DB.Prepare(updateScoreSQL)
	if err != nil {
		return newRating, fmt.Errorf("error preparing update statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	if err != nil {
		return newRating, fmt.Errorf("error executing update statement: %v", err)
	}

	newRatingSQL := fmt.Sprintf("SELECT * FROM rating WHERE UserID = %d", userID)
	rows, err := database.DB.Query(newRatingSQL)
	if err != nil {
		return newRating, fmt.Errorf("error executing new rating statement: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&newRating.ID, &newRating.UserID, &newRating.Score, &newRating.CreatedAt, &newRating.UpdatedAt)
		if err != nil {
			return newRating, fmt.Errorf("error scanning new rating row: %v", err)
		}
	}
	if increment {
		fmt.Println("User score incremented successfully")
	} else {
		fmt.Println("User score decremented successfully")
	}
	return newRating, nil
}
