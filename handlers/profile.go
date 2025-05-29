package handlers

import (
	"log"
	"net/http"

	"myproject/database"
	"myproject/middleware"
	"myproject/models"
	"myproject/utils"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		log.Println("Пользователь не аутентифицирован")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := middleware.GetCurrentUserID(r)
	log.Printf("Текущий ID пользователя: %d", userID)

	if userID == 0 {
		log.Println("Не удалось определить ID пользователя")
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Не удалось определить пользователя",
		})
		return
	}

	user, err := database.GetUserByID(userID)
	if err != nil {
		log.Printf("Ошибка загрузки данных пользователя: %v", err)
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Ошибка загрузки данных пользователя",
		})
		return
	}
	log.Printf("Загружен пользователь: %+v", user)

	posts, err := database.GetPostsByAuthor(userID)
	if err != nil {
		log.Printf("Ошибка загрузки постов: %v", err)
		posts = []models.Post{} // Возвращаем пустой массив вместо ошибки
	}
	log.Printf("Загружено %d постов пользователя", len(posts))

	data := map[string]interface{}{
		"User":            user,
		"Posts":           posts,
		"IsAuthenticated": true,
		"CurrentUserID":   userID,
	}

	utils.RenderTemplate(w, r, "profile.html", data)
}
