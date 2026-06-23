package utils


import (
	"net/http"
	"stage/models"
	"log"
)



func RenderError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	data := models.ErrorPage{
		Code:    code,
		Message: message,
	}
	if err := Tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		log.Println("template execution error:", err)
	}
}