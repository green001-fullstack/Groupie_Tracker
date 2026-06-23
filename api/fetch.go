package api

import (
	"stage/models"
	"sync"
)

type geoResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type ArtistCache struct{
	artistsCache []models.FullArtist
	artistById map[int]models.FullArtist
	mu sync.RWMutex
	geoCache *GeoCache
}

func NewArtistCache(geoCache *GeoCache) *ArtistCache{
	return &ArtistCache{
		artistById: make(map[int]models.FullArtist),
		geoCache: geoCache,
	}
}
