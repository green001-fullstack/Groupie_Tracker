package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"stage/service"
	"stage/models"
	"stage/utils"
	"strconv"
	"strings"
)

type SingleHandler struct {
	service *service.ArtistService
}

func NewSingleArtist(service *service.ArtistService) *SingleHandler {
	return &SingleHandler{
		service: service,
	}
}

func (h *SingleHandler) SingleArtistHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	newPath := strings.Trim(path, "/")
	pathSlice := strings.Split(newPath, "/")

	// artistsCache := h.service.GetAllArtists()

	if len(pathSlice) != 2 {
		utils.RenderError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	val := pathSlice[1]

	id, err := strconv.Atoi(val)
	if err != nil {
		utils.RenderError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	artist, found := h.service.GetArtistByID(id)
	if found {
		locJs, err := json.Marshal(artist.Locations)
		if err != nil {
			utils.RenderError(w, http.StatusInternalServerError, "JSON Marshal Error")
		}

		pageData := models.ArtistPageData{
			Artist:         artist.Artist,
			DatesLocations: artist.DatesLocations,
			Locations:      template.JS(string(locJs)),
		}

		err = utils.Tmpl.ExecuteTemplate(w, "artist.html", pageData)
		if err != nil {
			utils.RenderError(w, http.StatusInternalServerError, "Template Error")
			return
		}
	} else {
		utils.RenderError(w, http.StatusNotFound, "Artist not found")
	}
}
