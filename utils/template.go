package utils

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"myproject/middleware"
)

var (
	Templates *template.Template
)

func LoadTemplates() {
	funcMap := template.FuncMap{
		"formatDate": formatDate,
		"isAuthor":   isAuthor,
	}

	// Загруска шаблонов
	var err error
	Templates, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Ошибка загрузки шаблонов: %v", err)
	}
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}

	// Добавляем общие данные для всех шаблонов
	data["IsAuthenticated"] = middleware.IsAuthenticated(r)
	data["CurrentUsername"] = middleware.GetCurrentUsername(r)
	data["CurrentUserID"] = middleware.GetCurrentUserID(r)

	err := Templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Printf("Ошибка рендеринга шаблона %s: %v", tmpl, err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

// форматирует дату для отображения
func formatDate(t time.Time) string {
	return t.Local().Format("02.01.2006 15:04")
}

// проверка, является ли текущий пользователь автором
func isAuthor(currentUserID, authorID int) bool {
	return currentUserID == authorID
}
