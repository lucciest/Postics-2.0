package handlers

import (
	"net/http"
	"strconv"

	"myproject/database"
	"myproject/utils"

	"github.com/gorilla/mux"
)

// Панель администратора
func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetAllUsers()
	if err != nil {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Ошибка загрузки пользователей",
		})
		return
	}

	utils.RenderTemplate(w, r, "admin_dashboard.html", map[string]interface{}{
		"Users": users,
	})
}

// Блокировка пользователя
func BanUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	if err := database.BanUser(userID); err != nil {
		http.Error(w, "Ошибка блокировки пользователя", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// Разблокировка пользователя
func UnbanUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	if err := database.UnbanUser(userID); err != nil {
		http.Error(w, "Ошибка разблокировки пользователя", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// Удаление пользователя
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	if err := database.DeleteUser(userID); err != nil {
		http.Error(w, "Ошибка удаления пользователя", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetAllUsers()
	if err != nil {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Ошибка загрузки пользователей",
		})
		return
	}

	utils.RenderTemplate(w, r, "admin_users.html", map[string]interface{}{
		"Users": users,
	})
}

// список всех постов для админа
func AdminPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetAllPostsWithAuthors(0, 0) // 0,0 - без пагинации
	if err != nil {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Ошибка загрузки постов",
		})
		return
	}

	utils.RenderTemplate(w, r, "admin_posts.html", map[string]interface{}{
		"Posts": posts,
	})
}

// админское редактирование поста
func AdminEditPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	post, err := database.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		post.Title = r.FormValue("title")
		post.Anons = r.FormValue("anons")
		post.FullText = r.FormValue("full_text")

		if err := database.UpdatePost(post); err != nil {
			http.Error(w, "Ошибка обновления поста", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	author, _ := database.GetUserByID(int(post.AuthorID))

	utils.RenderTemplate(w, r, "admin_edit_post.html", map[string]interface{}{
		"Post":        post,
		"Author":      author,
		"IsAdminEdit": true, // Флаг что это админское редактирование
	})
}

// аааадминское удаление поста
func AdminDeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	if err := database.DeletePost(postID); err != nil {
		http.Error(w, "Ошибка удаления поста", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}
