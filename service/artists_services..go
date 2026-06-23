package service

import (
	"stage/api"
	"stage/models"
	"sort"
	"strconv"
	"strings"
)

type ArtistService struct{
	cache *api.ArtistCache
}

func NewArtistService( cache *api.ArtistCache) *ArtistService{
	return &ArtistService{
		cache : cache,
	}
}

func ( h *ArtistService) Sort(artists []models.FullArtist, sortType string) []models.FullArtist{

	switch sortType {
	case "ascending":
		sort.Slice(artists, func(i, j int) bool {
			return artists[i].Name < artists[j].Name
		})
	case "descending":
		sort.Slice(artists, func(i, j int) bool {
			return artists[i].Name > artists[j].Name
		})
	case "oldest":
		sort.Slice(artists, func(i, j int) bool {
			return artists[i].CreationDate < artists[j].CreationDate
		})
	case "newest":
		sort.Slice(artists, func(i, j int) bool {
			return artists[i].CreationDate > artists[j].CreationDate
		})
	}
	return artists
}

func (h *ArtistService) Paginate(artists []models.FullArtist, pageNo int, limit int)([]models.FullArtist, []int){
	if pageNo < 1 {
		pageNo= 1
	}

	totalPageNum := len(artists) / limit
	if len(artists)%limit != 0 {
		totalPageNum = totalPageNum + 1
	}
	start := (pageNo - 1) * limit
	if start > len(artists) {
		start = len(artists)
	}
	if start < 0 {
		start = 0
	}

	end := limit + start

	if end > len(artists) {
		end = len(artists)
	}

	pageSlice := artists[start:end]

	pageNumbers := make([]int, totalPageNum)

	for i := 0; i < totalPageNum; i++ {
		pageNumbers[i] = i + 1
	}
	return pageSlice, pageNumbers
}

func (h *ArtistService) Search( query string) []models.FullArtist{
	artistsCache := h.cache.GetAllArtists()
	
	query = strings.TrimSpace(query)
	
	var result []models.FullArtist
	if query != "" {
	query = strings.ToLower(query)

	for _, artist := range artistsCache {
		matched := false
		if strings.Contains(strings.ToLower(artist.Name), query) {
			matched = true
		}

		if !matched {
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), query) {
					matched = true
					break
				}
			}
		}

		if !matched {
			for location := range artist.DatesLocations {
				if strings.Contains(strings.ToLower(location), query) {
					matched = true
					break
				}
			}
		}

		if !matched {
			creationDate := strconv.Itoa(artist.CreationDate)
			if strings.Contains(creationDate, query) {
				matched = true
				// break
			}
		}

		if !matched {
			if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
				matched = true
				// break
			}
		}
		if matched {
			result = append(result, artist)
		}
	}
}else{
	result = artistsCache
}
	return result
}

func (s *ArtistService) GetArtists(query, sortType string, page, limit int) ([]models.FullArtist, []int) {

	artists := s.Search(query)
	artists = s.Sort(artists, sortType)

	return s.Paginate(artists, page, limit)
}

func (s *ArtistService) GetArtistByID(Id int) (models.FullArtist, bool){
	return s.cache.GetArtistByID(Id)
}