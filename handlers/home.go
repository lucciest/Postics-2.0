package handlers

import (
	"log"
	"net/http"

	"myproject/database"
	"myproject/middleware"
	"myproject/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetAllPosts(0, 0) // пока оставил 0, чтоб было без ограничений
	if err != nil {
		log.Printf("Ошибка при получении постов: %v", err)
		utils.RenderTemplate(w, r, "error.html", map[string]interface{}{
			"Error": "Ошибка при загрузке постов",
		})
		return
	}

	var postsWithAuthors []map[string]interface{}

	for _, post := range posts {
		author, err := database.GetUserByID(int(post.AuthorID))
		if err != nil {
			log.Printf("Ошибка получения автора: %v", err)
			continue
		}

		comments, err := database.GetCommentsByPostID(int(post.ID))
		commentCount := 0
		if err == nil {
			commentCount = len(comments)
		}

		postsWithAuthors = append(postsWithAuthors, map[string]interface{}{
			"Post":         post,
			"Author":       author,
			"CommentCount": commentCount,
		})
	}

	data := map[string]interface{}{
		"Posts":      postsWithAuthors,
		"User":       middleware.GetCurrentUser(r),
		"IsLoggedIn": middleware.IsAuthenticated(r),
	}

	utils.RenderTemplate(w, r, "index.html", data)
}
