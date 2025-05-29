package handlers

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"myproject/database"
	"myproject/middleware"
	"myproject/models"
	"myproject/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Проверка существования пользователя
		exists, err := database.CheckUserExists(username, email)
		if err != nil {
			log.Printf("Ошибка проверки пользователя: %v", err)
			utils.RenderTemplate(w, r, "register.html", map[string]interface{}{
				"Error": "Ошибка сервера",
			})
			return
		}
		if exists {
			utils.RenderTemplate(w, r, "register.html", map[string]interface{}{
				"Error": "Пользователь с таким именем или email уже существует",
			})
			return
		}

		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Ошибка хеширования пароля: %v", err)
			utils.RenderTemplate(w, r, "register.html", map[string]interface{}{
				"Error": "Ошибка сервера",
			})
			return
		}

		// Создание пользователя
		user := models.NewUser(username, string(hashedPassword), email)
		_, err = database.CreateUser(user)
		if err != nil {
			log.Printf("Ошибка создания пользователя: %v", err)
			utils.RenderTemplate(w, r, "register.html", map[string]interface{}{
				"Error": "Ошибка сервера",
			})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, r, "register.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := database.GetUserByUsername(username)
		if err != nil {
			log.Printf("Ошибка: %v", err)
			utils.RenderTemplate(w, r, "login.html", map[string]interface{}{
				"Error": "Неверные данные",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			utils.RenderTemplate(w, r, "login.html", map[string]interface{}{
				"Error": "Неверный пароль",
			})
			return
		}

		session, err := middleware.Store.Get(r, "session")
		if err != nil {
			log.Printf("Ошибка сессии: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		session.Values["authenticated"] = true
		session.Values["user_id"] = user.ID
		session.Values["username"] = user.Username

		if err := session.Save(r, w); err != nil {
			log.Printf("Ошибка сохранения сессии: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Printf("Успешный вход: %s", username) // Лог для отладки
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, r, "login.html", nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session")
	session.Values["authenticated"] = false
	session.Values["user_id"] = 0
	session.Values["username"] = ""
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
