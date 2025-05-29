package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	Store = sessions.NewCookieStore([]byte("secret-key-1"))
)

// проверяет аутентификацию пользователя
func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "session")
		if err != nil {
			log.Printf("Ошибка получения сессии: %v", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// проверяет, аутентифицирован ли пользователь
func IsAuthenticated(r *http.Request) bool {
	session, err := Store.Get(r, "session")
	if err != nil {
		log.Printf("Ошибка получения сессии: %v", err)
		return false
	}

	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

// возвращает имя текущего пользователя
func GetCurrentUsername(r *http.Request) string {
	session, err := Store.Get(r, "session")
	if err != nil {
		log.Printf("Ошибка получения сессии: %v", err)
		return ""
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		return ""
	}
	return username
}

// возвращает ID текущего пользователя
func GetCurrentUserID(r *http.Request) int {
	session, err := Store.Get(r, "session")
	if err != nil {
		log.Printf("Ошибка получения сессии: %v", err)
		return 0
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return 0
	}
	return userID
}

// GetCurrentUser возвращает данные текущего пользователя
func GetCurrentUser(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"ID":       GetCurrentUserID(r),
		"Username": GetCurrentUsername(r),
	}
}
