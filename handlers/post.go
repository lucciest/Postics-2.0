package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"myproject/database"
	"myproject/middleware"
	"myproject/models"
	"myproject/utils"
)

func ShowPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Неверный ID поста",
		})
		return
	}

	post, err := database.GetPostByID(postID)
	if err != nil {
		log.Printf("Ошибка получения поста: %v", err)
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Пост не найден",
		})
		return
	}

	author, err := database.GetUserByID(int(post.AuthorID))
	if err != nil {
		log.Printf("Ошибка получения автора: %v", err)
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Ошибка загрузки автора",
		})
		return
	}

	comments, err := database.GetCommentsByPostID(postID)
	if err != nil {
		log.Printf("Ошибка получения комментариев: %v", err)
		comments = []models.Comment{} // Пустой массив вместо ошибки
	}

	data := map[string]interface{}{
		"Post":            post,
		"Author":          author,
		"Comments":        comments,
		"CurrentUser":     middleware.GetCurrentUser(r),
		"IsAuthenticated": middleware.IsAuthenticated(r),
	}

	utils.RenderTemplate(w, r, "post.html", data)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		title := r.FormValue("title")
		anons := r.FormValue("anons")
		fullText := r.FormValue("full_text")

		if title == "" || anons == "" || fullText == "" {
			utils.RenderTemplate(w, r, "create_post.html", map[string]interface{}{
				"Error": "Все поля должны быть заполнены",
			})
			return
		}

		userID := middleware.GetCurrentUserID(r)
		post := &models.Post{
			Title:    title,
			Anons:    anons,
			FullText: fullText,
			AuthorID: userID,
		}

		// Исправленная строка - получаем оба возвращаемых значения
		id, err := database.CreatePost(post)
		if err != nil {
			log.Printf("Ошибка создания поста: %v", err)
			utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
				"Error": "Ошибка при создании поста",
			})
			return
		}

		http.Redirect(w, r, "/post/"+strconv.Itoa(int(id)), http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, r, "create_post.html", nil)
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Неверный ID поста",
		})
		return
	}

	post, err := database.GetPostByID(postID)
	if err != nil {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Пост не найден",
		})
		return
	}

	currentUserID := middleware.GetCurrentUserID(r)
	if post.AuthorID != currentUserID {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "У вас нет прав для редактирования этого поста",
		})
		return
	}

	if r.Method == "POST" {
		title := r.FormValue("title")
		anons := r.FormValue("anons")
		fullText := r.FormValue("full_text")

		post.Title = title
		post.Anons = anons
		post.FullText = fullText

		if err := database.UpdatePost(post); err != nil {
			log.Printf("Ошибка обновления поста: %v", err)
			utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
				"Error": "Ошибка при обновлении поста",
			})
			return
		}

		http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Post": post,
	}

	utils.RenderTemplate(w, r, "edit_post.html", data)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Неверный ID поста",
		})
		return
	}

	post, err := database.GetPostByID(postID)
	if err != nil {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Пост не найден",
		})
		return
	}

	currentUserID := middleware.GetCurrentUserID(r)
	if post.AuthorID != currentUserID {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "У вас нет прав для удаления этого поста",
		})
		return
	}

	if err := database.DeletePost(postID); err != nil {
		log.Printf("Ошибка удаления поста: %v", err)
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Ошибка при удалении поста",
		})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
