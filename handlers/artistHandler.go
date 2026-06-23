package handlers

import (
	"net/http"
	"stage/models"
	"stage/utils"
	"strconv"
	"stage/service"
)

type Handler struct{
	Service *service.ArtistService
}

func(h *Handler) ArtistsHandler(w http.ResponseWriter, r *http.Request){

	query := r.URL.Query().Get("query")
	sortType := r.URL.Query().Get("sort")
	pageNum, _ := strconv.Atoi(r.URL.Query().Get("page"))

	artists, pageNo := h.Service.GetArtists(query, sortType, pageNum, 16)

	data := models.ArtistsPageData{
		Artists:      artists,
		PageNumbers:  pageNo,
		PageNo:       pageNum,
		Query:        query,
		Sort:         sortType,
		TotalArtists: len(artists),
	}

	err := utils.Tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		utils.RenderError(w, http.StatusInternalServerError, "Template Error")
		return
	}
}

	
