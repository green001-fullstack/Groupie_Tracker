package handlers


import (
	"net/http"
    "stage/utils"
)


func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		utils.RenderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if r.URL.Path != "/" {
		utils.RenderError(w, http.StatusNotFound, "Page Not found")
		return
	}
	err := utils.Tmpl.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		utils.RenderError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}