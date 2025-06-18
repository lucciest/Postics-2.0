package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"myproject/config"
	"myproject/database"
	"myproject/handlers"
	"myproject/middleware"
	"myproject/utils"
)

func main() {
	// Инициализация конфигурации
	config.LoadConfig()

	// Подключение к базе данных
	if err := database.InitDB(config.AppConfig.DatabaseURL); err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer database.CloseDB()

	// Загрузка HTML-шаблонов
	utils.LoadTemplates()

	// Настройка сессий
	middleware.Store = sessions.NewCookieStore([]byte(config.AppConfig.SessionKey))
	middleware.Store.Options = &sessions.Options{
		MaxAge:   86400 * 7, // 1 неделя
		HttpOnly: true,
		Path:     "/",
	}

	// Создание роутера
	r := mux.NewRouter()

	// Глобальные middleware
	r.Use(middleware.LoggingMiddleware)

	// Статические файлы (CSS, JS)
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	)

	// Публичные маршруты
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/home", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/post/{id:[0-9]+}", handlers.ShowPostHandler).Methods("GET")

	// Аутентификация
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")

	// Проверка аутентификации внутри обработчиков
	r.HandleFunc("/create-post", handlers.CreatePostHandler).Methods("GET", "POST")
	r.HandleFunc("/edit-post/{id:[0-9]+}", handlers.EditPostHandler).Methods("GET", "POST")
	r.HandleFunc("/delete-post/{id:[0-9]+}", handlers.DeletePostHandler).Methods("POST")
	r.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")

	// Комментарии
	r.HandleFunc("/post/{id:[0-9]+}/comment", handlers.AddCommentHandler).Methods("POST")
	r.HandleFunc("/comment/{id:[0-9]+}/delete", handlers.DeleteCommentHandler).Methods("POST")
	r.HandleFunc("/comment/{id:[0-9]+}/edit", handlers.EditCommentHandler).Methods("GET", "POST")

	// Админские маршруты
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthRequired)
	adminRouter.Use(middleware.AdminRequired)
	// работа с пользователями
	adminRouter.HandleFunc("", handlers.AdminDashboardHandler).Methods("GET")
	adminRouter.HandleFunc("/users", handlers.AdminUsersHandler).Methods("GET") // Новый обработчик
	adminRouter.HandleFunc("/ban", handlers.BanUserHandler).Methods("POST")
	adminRouter.HandleFunc("/unban", handlers.UnbanUserHandler).Methods("POST")
	adminRouter.HandleFunc("/delete-user", handlers.DeleteUserHandler).Methods("POST")
	// работа с постами
	adminRouter.HandleFunc("/posts", handlers.AdminPostsHandler).Methods("GET")
	adminRouter.HandleFunc("/posts/{id:[0-9]+}/edit", handlers.AdminEditPostHandler).Methods("GET", "POST")
	adminRouter.HandleFunc("/posts/delete", handlers.AdminDeletePostHandler).Methods("POST")

	// Запуск сервера
	addr := ":" + config.AppConfig.ServerPort
	log.Printf("Сервер запущен на http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
