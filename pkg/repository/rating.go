package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// Import your models if needed
	"PoliticianRating/pkg/model"
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

func (r *RatingRepository) IncrementRating(userId int) {
	incrementScoreSQL := `UPDATE rating SET Score = Score + 1, UpdatedAt = ? WHERE UserID = ?`
	stmt, err := r.Db.Prepare(incrementScoreSQL)
	if err != nil {
		fmt.Println("Error preparing update statement:", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(time.Now(), userId)
	if err != nil {
		fmt.Println("Error executing update statement:", err)
		return
	}

	fmt.Println("User score incremented successfully")
}

func (r *RatingRepository) DecrementRating(userId int) {
	incrementScoreSQL := `UPDATE rating SET Score = Score - 1 WHERE UserID = ?`
	stmt, err := r.Db.Prepare(incrementScoreSQL)
	if err != nil {
		fmt.Println("Error preparing update statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		fmt.Println("Error executing update statement:", err)
		return
	}

	fmt.Println("User score decremented successfully")
}

func (r *RatingRepository) GetAllRatings(order string) ([]model.Rating, error) {
	var ratings []model.Rating
	var orderBy string
	switch order {
	case "asc":
		orderBy = "ORDER BY Score ASC"
	case "desc":
		orderBy = "ORDER BY Score DESC"
	default:
		orderBy = "" // Default to no sorting
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

func (r *RatingRepository) DeleteUser(userID int) error {
	deleteSQL := `DELETE FROM rating WHERE UserID = ?`

	_, err := r.Db.Exec(deleteSQL, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
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

	//for row.Next() {
	err2 := row.Scan(&rating.ID, &rating.UserID, &rating.Score, &rating.CreatedAt, &rating.UpdatedAt)
	if err2 != nil {
		return rating, fmt.Errorf("error scanning ratings: %v", err2)
	}
	//}
	return rating, nil
}

func (r *RatingRepository) UpdateUserRating(userID int, increment bool) (model.Rating, error) {
	var newRating model.Rating

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
