package database

import (
	"myproject/models"
	"time"
)

func GetCommentByID(id int) (*models.Comment, error) {
	var comment models.Comment
	err := DB.QueryRow(`
		SELECT id, post_id, user_id, username, content, created_at 
		FROM comments WHERE id = ?
	`, id).Scan(
		&comment.ID, &comment.PostID, &comment.UserID, &comment.Username,
		&comment.Content, &comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
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
		SET content = ?
		WHERE id = ?
	`, comment.Content, comment.ID)
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
