package database

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./foo.db?parseTime=true")
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

	_, err = DB.Exec(createTableSQL)
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

	_, err = DB.Exec(createTriggerSQL)
	if err != nil {
		log.Fatalf("Failed to create trigger: %v", err)
	}
	//defer db.Close() 			???
	fmt.Println("Database setup completed")
	return nil
}
