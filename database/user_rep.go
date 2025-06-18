package database

import (
	"time"

	"myproject/models"
)

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	var createdAt []byte

	err := DB.QueryRow(`
        SELECT id, username, password, email, is_admin, created_at 
        FROM users WHERE id = ?
    `, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email,
		&user.IsAdmin, &createdAt,
	)
	if err != nil {
		return nil, err
	}

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	var createdAt []byte // Временная переменная для хранения даты
	var isBanned int

	err := DB.QueryRow(`
        SELECT id, username, password, email, is_admin, is_banned, created_at 
        FROM users WHERE username = ?
    `, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email,
		&user.IsAdmin, &isBanned, &createdAt,
	)
	if err != nil {
		return nil, err
	}

	// Конвертируем int в bool
	user.IsBanned = isBanned != 0

	// Парсим дату из байтов
	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(user *models.User) (int, error) {
	res, err := DB.Exec(`
		INSERT INTO users (username, password, email, created_at)
		VALUES (?, ?, ?, ?)
	`, user.Username, user.Password, user.Email, time.Now())

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func UpdateUser(user *models.User) error {
	_, err := DB.Exec(`
		UPDATE users
		SET username = ?, password = ?, email = ?
		WHERE id = ?
	`, user.Username, user.Password, user.Email, user.ID)
	return err
}

func DeleteUser(id int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Удаление комментариев пользователя
	_, err = tx.Exec("DELETE FROM comments WHERE user_id = ?", id)
	if err != nil {
		return err
	}

	// Удаление комментариев к постам пользователя
	_, err = tx.Exec("DELETE FROM comments WHERE post_id IN (SELECT id FROM articles WHERE author_id = ?)", id)
	if err != nil {
		return err
	}

	// Удаление постов пользователя
	_, err = tx.Exec("DELETE FROM articles WHERE author_id = ?", id)
	if err != nil {
		return err
	}

	// Удаление пользователя
	_, err = tx.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// проверка существования пользователя
func CheckUserExists(username, email string) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) 
		FROM users 
		WHERE username = ? OR email = ?
	`, username, email).Scan(&count)
	return count > 0, err
}
