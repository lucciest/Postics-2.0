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

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Для добавления комментария необходимо войти в систему",
		})
		return
	}

	if r.Method != "POST" {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Метод не поддерживается",
		})
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

	content := r.FormValue("content")
	if content == "" {
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Комментарий не может быть пустым",
		})
		return
	}

	userID := middleware.GetCurrentUserID(r)
	username := middleware.GetCurrentUsername(r)

	comment := models.NewComment(postID, userID, username, content)
	_, err = database.CreateComment(comment)
	if err != nil {
		log.Printf("Ошибка при создании комментария: %v", err)
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Ошибка при добавлении комментария",
		})
		return
	}

	http.Redirect(w, r, "/post/"+vars["id"], http.StatusSeeOther)
}

func EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Error(w, "Для редактирования комментария необходимо войти в систему", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		http.Error(w, "Комментарий не найден", http.StatusNotFound)
		return
	}

	currentUserID := middleware.GetCurrentUserID(r)
	if comment.UserID != currentUserID {
		http.Error(w, "У вас нет прав для редактирования этого комментария", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		content := r.FormValue("content")
		if content == "" {
			http.Error(w, "Комментарий не может быть пустым", http.StatusBadRequest)
			return
		}

		comment.Content = content
		err = database.UpdateComment(comment)
		if err != nil {
			log.Printf("Ошибка при обновлении комментария: %v", err)
			http.Error(w, "Ошибка при редактировании комментария", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post/"+strconv.Itoa(comment.PostID), http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, r, "edit_comment.html", map[string]interface{}{
		"Comment": comment,
	})
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Error(w, "Для удаления комментария необходимо войти в систему", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		http.Error(w, "Комментарий не найден", http.StatusNotFound)
		return
	}

	currentUserID := middleware.GetCurrentUserID(r)
	if comment.UserID != currentUserID {
		http.Error(w, "У вас нет прав для удаления этого комментария", http.StatusForbidden)
		return
	}

	err = database.DeleteComment(commentID)
	if err != nil {
		log.Printf("Ошибка при удалении комментария: %v", err)
		http.Error(w, "Ошибка при удалении комментария", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(comment.PostID), http.StatusSeeOther)
}
