package database

import (
	"database/sql"
	"myproject/models"
	"time"
)

func GetPostByID(id int) (*models.Post, error) {
	var post models.Post
	var fullText sql.NullString // Используем NullString для обработки NULL значений
	var createdAt []byte

	err := DB.QueryRow(`
        SELECT id, title, anons, full_text, author_id, created_at 
        FROM articles WHERE id = ?
    `, id).Scan(
		&post.ID, &post.Title, &post.Anons, &fullText,
		&post.AuthorID, &createdAt,
	)
	if err != nil {
		return nil, err
	}

	// Обрабатываем возможный NULL в full_text
	if fullText.Valid {
		post.FullText = fullText.String
	} else {
		post.FullText = ""
	}

	post.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func CreatePost(post *models.Post) (int, error) {
	res, err := DB.Exec(`
		INSERT INTO articles (title, anons, full_text, author_id, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, post.Title, post.Anons, post.FullText, post.AuthorID, time.Now())

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func UpdatePost(post *models.Post) error {
	_, err := DB.Exec(`
		UPDATE articles
		SET title = ?, anons = ?, full_text = ?
		WHERE id = ?
	`, post.Title, post.Anons, post.FullText, post.ID)
	return err
}

func DeletePost(id int) error {
	// Удаляем сначала комментарии к посту
	if err := DeleteCommentsByPostID(id); err != nil {
		return err
	}

	// Затем удаляем сам пост
	_, err := DB.Exec("DELETE FROM articles WHERE id = ?", id)
	return err
}

func GetAllPosts(limit, offset int) ([]models.Post, error) {
	query := `
        SELECT id, title, anons, full_text, author_id, created_at
        FROM articles
        ORDER BY created_at DESC
    `

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		var createdAt []byte

		err := rows.Scan(
			&p.ID, &p.Title, &p.Anons, &p.FullText,
			&p.AuthorID, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		p.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func GetPostsByAuthor(authorID int) ([]models.Post, error) {
	rows, err := DB.Query(`
        SELECT id, title, anons, full_text, created_at
        FROM articles
        WHERE author_id = ?
        ORDER BY created_at DESC
    `, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		var createdAt []byte // Временная переменная для хранения даты

		err := rows.Scan(
			&p.ID, &p.Title, &p.Anons, &p.FullText,
			&createdAt, // Сканируем во временную переменную
		)
		if err != nil {
			return nil, err
		}

		// Парсим дату из байтов
		p.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, err
		}
		p.AuthorID = authorID

		posts = append(posts, p)
	}

	return posts, nil
}

func GetAllPostsWithAuthors(limit, offset int) ([]struct {
	models.Post
	AuthorUsername string
}, error) {
	query := `
        SELECT a.id, a.title, a.anons, a.full_text, a.author_id, 
               a.created_at, u.username as author_username
        FROM articles a
        JOIN users u ON a.author_id = u.id
        ORDER BY a.created_at DESC
    `

	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
	}

	var rows *sql.Rows
	var err error

	if limit > 0 {
		rows, err = DB.Query(query, limit, offset)
	} else {
		rows, err = DB.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []struct {
		models.Post
		AuthorUsername string
	}

	for rows.Next() {
		var p struct {
			models.Post
			AuthorUsername string
		}
		var createdAt []byte

		err := rows.Scan(
			&p.ID, &p.Title, &p.Anons, &p.FullText,
			&p.AuthorID, &createdAt, &p.AuthorUsername,
		)
		if err != nil {
			return nil, err
		}

		p.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}
