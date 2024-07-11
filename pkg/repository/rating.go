package repository

import (
	"PoliticianRating/pkg/model"
	"database/sql"
	"fmt"
	"log"
)

type RatingRepository struct {
	Db *sql.DB
}

func NewRatingRepository() *RatingRepository {
	db, err := sql.Open("sqlite3", "./foo.db?parseTime=true")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS rating (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		UserID INTEGER,
		Score INTEGER,
		CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	createTriggerSQL := `
	CREATE TRIGGER IF NOT EXISTS update_timestamp
	AFTER UPDATE ON rating
	FOR EACH ROW
	BEGIN
		UPDATE rating
		SET UpdatedAt = CURRENT_TIMESTAMP
		WHERE id = OLD.id;
	END;`

	_, err = db.Exec(createTriggerSQL)
	if err != nil {
		log.Fatalf("Failed to create trigger: %v", err)
	}

	//defer db.Close() 			???
	fmt.Println("Database setup completed")
	return &RatingRepository{Db: db}
}

func (r *RatingRepository) GetAllRatings(order string) ([]model.Rating, error) {
	var ratings []model.Rating
	var orderBy string
	switch order {
	case "+score":
		orderBy = "ORDER BY Score ASC"
	case "-score":
		orderBy = "ORDER BY Score DESC"
	case "+date":
		orderBy = "ORDER BY UpdatedAt ASC"
	case "-date":
		orderBy = "ORDER BY UpdatedAt DESC"
	default:
		return nil, fmt.Errorf("invalid order: %s", order)
	}

	selectAllSQL := fmt.Sprintf("SELECT ID, UserId, Score, CreatedAt, UpdatedAt FROM rating %s", orderBy)
	rows, err := r.Db.Query(selectAllSQL)
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

func (r *RatingRepository) GetUserRating(userID int) (model.Rating, error) {
	var rating model.Rating
	selectUserSql := fmt.Sprintf("SELECT * FROM rating WHERE UserID = %d", userID)
	row, err := r.Db.Query(selectUserSql)
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
func (r *RatingRepository) UpdateUserRating(userID int, increment bool) (model.Rating, error) {
	var newRating model.Rating

	var userExists bool
	err := r.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM rating WHERE UserID = ?)", userID).Scan(&userExists)
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

	stmt, err := r.Db.Prepare(updateScoreSQL)
	if err != nil {
		return newRating, fmt.Errorf("error preparing update statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	if err != nil {
		return newRating, fmt.Errorf("error executing update statement: %v", err)
	}

	newRatingSQL := fmt.Sprintf("SELECT * FROM rating WHERE UserID = %d", userID)
	rows, err := r.Db.Query(newRatingSQL)
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

func (r *RatingRepository) InsertRating(userID int, score int) error {
	var exists bool
	checkUserSQL := `SELECT EXISTS(SELECT 1 FROM rating WHERE userId = ?)`
	err := r.Db.QueryRow(checkUserSQL, userID).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking for user:", err)
		return err
	}
	newRating := model.Rating{
		UserID: userID,
		Score:  score,
	}
	if !exists {
		insertRatingSQL := `INSERT INTO rating (UserID, Score) VALUES (?, ?)`
		_, err := r.Db.Exec(insertRatingSQL, newRating.UserID, newRating.Score)
		if err != nil {
			return fmt.Errorf("error inserting rating: %v", err)
		}
		return nil
	}
	return fmt.Errorf("user with ID %d already has a rating", userID)
}

func (r *RatingRepository) DeleteUser(userID int) error {
	deleteSQL := `DELETE FROM rating WHERE UserID = ?`

	_, err := r.Db.Exec(deleteSQL, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}
