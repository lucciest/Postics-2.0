package database

import (
	"log"
	"myproject/models"
	"time"
)

func GetCommentByID(id int) (*models.Comment, error) {
	log.Printf("Поиск комментария с ID: %d", id)

	var comment models.Comment
	var createdAt []byte

	err := DB.QueryRow(`
        SELECT id, post_id, user_id, username, content, created_at 
        FROM comments WHERE id = ?
    `, id).Scan(
		&comment.ID, &comment.PostID, &comment.UserID, &comment.Username,
		&comment.Content, &createdAt,
	)

	if err != nil {
		log.Printf("Ошибка при поиске комментария: %v", err)
		return nil, err
	}

	// парсим
	parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		log.Printf("Ошибка парсинга даты: %v", err)
		return nil, err
	}
	comment.CreatedAt = parsedTime

	log.Printf("Комментарий найден: %+v", comment)
	return &comment, nil
}

func CreateComment(comment *models.Comment) (int, error) {
	res, err := DB.Exec(`
		INSERT INTO comments (post_id, user_id, username, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, comment.PostID, comment.UserID, comment.Username, comment.Content, time.Now())

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func UpdateComment(comment *models.Comment) error {
	_, err := DB.Exec(`
        UPDATE comments
        SET content = ?, username = ?
        WHERE id = ?
    `, comment.Content, comment.Username, comment.ID)
	return err
}

func DeleteComment(id int) error {
	_, err := DB.Exec("DELETE FROM comments WHERE id = ?", id)
	return err
}

func GetCommentsByPostID(postID int) ([]models.Comment, error) {
	rows, err := DB.Query(`
        SELECT id, post_id, user_id, username, content, created_at
        FROM comments
        WHERE post_id = ?
        ORDER BY created_at ASC
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		var createdAt []byte

		err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Username,
			&c.Content, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		// Парсим дату из байтов
		c.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}

func GetCommentsByUserID(userID int) ([]models.Comment, error) {
	rows, err := DB.Query(`
		SELECT id, post_id, user_id, username, content, created_at
		FROM comments
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Username,
			&c.Content, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func DeleteCommentsByPostID(postID int) error {
	_, err := DB.Exec("DELETE FROM comments WHERE post_id = ?", postID)
	return err
}
