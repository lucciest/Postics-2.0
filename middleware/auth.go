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
			log.Printf("Ошибка сессии: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// проверка блокировки через сессию
		if banned, ok := session.Values["is_banned"].(bool); ok && banned {
			session.Values["authenticated"] = false
			session.Options.MaxAge = -1
			session.Save(r, w)
			http.Redirect(w, r, "/login?banned=1", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "session")
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		isAdmin, ok := session.Values["is_admin"].(bool)
		if !ok || !isAdmin {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
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

// возвращает данные текущего пользователя
func GetCurrentUser(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"ID":       GetCurrentUserID(r),
		"Username": GetCurrentUsername(r),
	}
}

func IsAdmin(r *http.Request) bool {
	session, err := Store.Get(r, "session")
	if err != nil {
		log.Printf("Session error: %v", err)
		return false
	}

	isAdmin, ok := session.Values["is_admin"].(bool)
	if !ok {
		// Для случаев, когда значение хранится как int
		if adminInt, ok := session.Values["is_admin"].(int64); ok {
			return adminInt == 1
		}
		return false
	}
	return isAdmin
}
