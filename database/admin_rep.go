package database

import (
	"myproject/models"
	"time"
)

// Блокировка пользователя
func BanUser(userID int) error {
	_, err := DB.Exec("UPDATE users SET is_banned = TRUE WHERE id = ?", userID)
	return err
}

// Разблокировка пользователя
func UnbanUser(userID int) error {
	_, err := DB.Exec("UPDATE users SET is_banned = FALSE WHERE id = ?", userID)
	return err
}

// Получение списка всех пользователей (для админ-панели)
func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query(`
        SELECT id, username, email, is_admin, is_banned, created_at 
        FROM users
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var createdAt []byte

		err := rows.Scan(
			&user.ID, &user.Username, &user.Email,
			&user.IsAdmin, &user.IsBanned, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
